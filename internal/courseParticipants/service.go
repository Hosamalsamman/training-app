package courseParticipants

import (
	"errors"
	"training-app/internal/courses"
	"training-app/internal/dbutil"
	"training-app/internal/models"
	"training-app/internal/pathGradeSubjects"
	"training-app/internal/persons"

	"gorm.io/gorm"
)

// Sentinel errors for the participant rules. The handler maps
// them to proper HTTP status codes.
var (
	// The referenced course does not exist for this client.
	ErrCourseNotFound = errors.New("course not found for this client")

	// The referenced person does not exist for this client.
	ErrPersonNotFound = errors.New("person not found for this client")

	// The participant id does not exist for this client.
	ErrParticipantNotFound = errors.New("course participant not found for this client")

	// The person is not an eligible path grade candidate for
	// the course: no candidate row for the same path grade
	// (and term), or the candidate already passed.
	ErrNotEligiblePathGradeCandidate = errors.New("the person is not an eligible path grade candidate for this course")
)

// The candidate rule no longer uses hardcoded evaluation
// ids: the client-owned evaluation_types table carries a
// should_repeat column, so eligibility is data-driven per
// client. should_repeat = 0 (passed: ناجح, ناجح بامتياز)
// blocks the course; should_repeat = 1 (failed) or NULL
// allows it.

type Service struct {
	repo *Repository

	// The courses and persons modules own their tables: the
	// course and the person a participant references are
	// loaded through their repositories instead of this
	// module re-implementing those lookups.
	courseRepo *courses.Repository
	personRepo *persons.Repository

	// The pathGradeSubjects module owns the subject table,
	// used to resolve the path grade of the course's subject
	// inside the candidate rule.
	pathGradeSubjectRepo *pathGradeSubjects.Repository
}

func NewService(repo *Repository, courseRepo *courses.Repository, personRepo *persons.Repository, pathGradeSubjectRepo *pathGradeSubjects.Repository) *Service {
	return &Service{
		repo:                 repo,
		courseRepo:           courseRepo,
		personRepo:           personRepo,
		pathGradeSubjectRepo: pathGradeSubjectRepo,
	}
}

func (s *Service) GetAll(clientID int) ([]models.CourseParticipant, error) {

	repo := s.repo.ForClient(clientID)

	return repo.GetAll()
}

func (s *Service) GetByID(clientID int, id int) (*models.CourseParticipant, error) {

	repo := s.repo.ForClient(clientID)

	return repo.GetByID(id)
}

// Create validates the request, applies the path-grade
// candidate rule against the course the participant belongs
// to, and inserts it, returning it with all relations
// preloaded.
func (s *Service) Create(clientID int, req *models.CourseParticipantRequest) (*models.CourseParticipant, error) {

	var participant models.CourseParticipant

	err := dbutil.WithTransaction(
		s.repo.DB(),
		func(tx *gorm.DB) error {

			// Use the transaction DB and restrict this
			// repository to the authenticated client's data.
			repo := s.repo.
				WithDB(tx).
				ForClient(clientID)

			course, err := s.loadCourseForValidation(tx, clientID, req.CourseID)

			if err != nil {
				return err
			}

			person, err := s.loadPersonForValidation(tx, clientID, req.PersonID)

			if err != nil {
				return err
			}

			// Rule: a path-grade course can only be taken
			// by an eligible candidate of the same path
			// grade (and term). The check runs before the
			// INSERT so an invalid person aborts the whole
			// transaction.
			if err := s.validatePathGradeCandidate(tx, clientID, person.ID, course); err != nil {
				return err
			}

			participant = models.CourseParticipant{
				CourseID: req.CourseID,
				PersonID: person.ID,
				ClientID: clientID,
			}

			if err := repo.Create(&participant); err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return s.repo.ForClient(clientID).GetByID(participant.ID)
}

// Update re-applies the same candidate rule of creation to the
// new values, so a participant can never be moved to a course
// they are not eligible for.
func (s *Service) Update(clientID int, id int, req *models.CourseParticipantRequest) (*models.CourseParticipant, error) {

	err := dbutil.WithTransaction(
		s.repo.DB(),
		func(tx *gorm.DB) error {

			repo := s.repo.
				WithDB(tx).
				ForClient(clientID)

			course, err := s.loadCourseForValidation(tx, clientID, req.CourseID)

			if err != nil {
				return err
			}

			person, err := s.loadPersonForValidation(tx, clientID, req.PersonID)

			if err != nil {
				return err
			}

			if err := s.validatePathGradeCandidate(tx, clientID, person.ID, course); err != nil {
				return err
			}

			participant := models.CourseParticipant{
				CourseID: req.CourseID,
				PersonID: person.ID,
			}

			rows, err := repo.Update(id, &participant)

			if err != nil {
				return err
			}

			if rows == 0 {
				return ErrParticipantNotFound
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return s.repo.ForClient(clientID).GetByID(id)
}

// Delete removes a participant by id. The client-scoped
// repository makes the id of another client's participant
// behave exactly like a missing one, and RowsAffected catches
// both cases.
func (s *Service) Delete(clientID int, id int) error {

	return dbutil.WithTransaction(
		s.repo.DB(),
		func(tx *gorm.DB) error {

			repo := s.repo.
				WithDB(tx).
				ForClient(clientID)

			rows, err := repo.Delete(id)

			if err != nil {
				return err
			}

			if rows == 0 {
				return ErrParticipantNotFound
			}

			return nil
		},
	)
}

// loadCourseForValidation loads the referenced course through
// the courses module repository scoped to the same transaction
// and client. The courses repository preloads the subject
// relations, but the rule only needs the subject ids of the
// row itself.
func (s *Service) loadCourseForValidation(tx *gorm.DB, clientID int, courseID int) (*models.Course, error) {

	course, err := s.courseRepo.
		WithDB(tx).
		ForClient(clientID).
		GetByID(courseID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourseNotFound
		}
		return nil, err
	}

	return course, nil
}

// loadPersonForValidation loads the referenced person through
// the persons module repository scoped to the same transaction
// and client.
func (s *Service) loadPersonForValidation(tx *gorm.DB, clientID int, personID int) (*models.Person, error) {

	person, err := s.personRepo.
		WithDB(tx).
		ForClient(clientID).
		GetByID(personID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPersonNotFound
		}
		return nil, err
	}

	return person, nil
}

// validatePathGradeCandidate enforces the candidate rule:
//
//   - a path_grade_subject course requires a candidate row of
//     the person for the subject's path grade;
//   - a path_grade_subject_term course additionally requires
//     the candidate's term to equal the subject term's term;
//   - the matching candidate's final evaluation must be NULL
//     or an evaluation whose should_repeat = 1 (failed), so a
//     person who already passed (should_repeat = 0) cannot
//     attend the course again.
//
// A learning_subject course carries no path grade, so no
// candidate rule applies to it. The check is a single counted
// query inside the transaction: "no candidate" and "candidate
// already passed" collapse into the same business failure.
func (s *Service) validatePathGradeCandidate(tx *gorm.DB, clientID int, personID int, course *models.Course) error {

	// A course that is neither a path-grade subject course
	// nor a subject-term course carries no path grade, so
	// anyone may attend it.
	if course.PathGradeSubjectID == nil && course.PathGradeSubjectTermID == nil {
		return nil
	}

	q := tx.Model(&models.PathGradeCandidate{}).
		Where("path_grade_candidates.person_id = ?", personID).
		Where("path_grade_candidates.client_id = ?", clientID).
		// NULL or a non-passing evaluation: either the
		// person has not been evaluated yet, or they took
		// the course and failed (should_repeat = 1), so
		// they may take it again. A passing evaluation
		// (should_repeat = 0, e.g. ناجح, ناجح بامتياز)
		// blocks the course. The should_repeat flag lives
		// on the client-owned evaluation_types row, so no
		// evaluation ids are hardcoded.
		Where("path_grade_candidates.final_evaluation IS NULL OR NOT EXISTS (SELECT 1 FROM evaluation_types AS et WHERE et.id = path_grade_candidates.final_evaluation AND et.should_repeat = ?)", false)

	if course.PathGradeSubjectID != nil {
		// Resolve the path grade of the subject through
		// the pathGradeSubjects module, client-scoped, so
		// a course pointing at another client's subject
		// fails the lookup as a missing course.
		subject, err := s.pathGradeSubjectRepo.
			WithDB(tx).
			ForClient(clientID).
			GetByIDPlain(*course.PathGradeSubjectID)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrCourseNotFound
			}
			return err
		}

		q = q.Where("path_grade_candidates.path_grade_id = ?", subject.PathGradeID)
	}

	if course.PathGradeSubjectTermID != nil {
		// The subject-term row carries both its subject
		// (and therefore the path grade) and the term the
		// course belongs to. Resolve both from the row,
		// with the subject fetched through the same
		// client-scoped module.
		term, err := s.loadSubjectTermForValidation(tx, clientID, *course.PathGradeSubjectTermID)

		if err != nil {
			return err
		}

		subject, err := s.pathGradeSubjectRepo.
			WithDB(tx).
			ForClient(clientID).
			GetByIDPlain(term.PathGradeSubjectID)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrCourseNotFound
			}
			return err
		}

		q = q.Where("path_grade_candidates.path_grade_id = ?", subject.PathGradeID)
		q = q.Where("path_grade_candidates.term_id = ?", term.TermID)
	}

	var count int64

	if err := q.Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return ErrNotEligiblePathGradeCandidate
	}

	return nil
}

// loadSubjectTermForValidation fetches the path grade subject
// term row of the course inside the transaction. Only the
// foreign keys of the row (subject id and term id) are needed,
// and the client_id filter keeps other clients' terms out.
func (s *Service) loadSubjectTermForValidation(tx *gorm.DB, clientID int, subjectTermID int) (*models.PathGradeSubjectTerm, error) {

	var term models.PathGradeSubjectTerm

	err := tx.Where("id = ? AND client_id = ?", subjectTermID, clientID).
		First(&term).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourseNotFound
		}
		return nil, err
	}

	return &term, nil
}