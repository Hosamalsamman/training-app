package pathGradeCandidate

import (
	"errors"
	"time"
	"training-app/internal/dbutil"
	"training-app/internal/models"
	"training-app/internal/pathGrades"
	"training-app/internal/persons"

	"gorm.io/gorm"
)

// Sentinel errors for the candidate rules. The handler maps
// them to proper HTTP status codes.
var (
	// The candidate id does not exist for this client.
	ErrPathGradeCandidateNotFound = errors.New("path grade candidate not found for this client")

	// The referenced person does not exist for this client.
	ErrPersonNotFound = errors.New("person not found for this client")

	// The referenced path grade does not exist for this client.
	ErrPathGradeNotFound = errors.New("path grade not found for this client")

	// The person has no current grade, so the next-grade
	// and interval rules cannot be evaluated.
	ErrPersonHasNoCurrentGrade = errors.New("person has no current grade")

	// The person has no date for the current grade, so the
	// interval rule cannot be evaluated.
	ErrPersonHasNoGradeDate = errors.New("person has no date for the current grade")

	// The current or target grade row is missing its numeric
	// code, so the next-grade rule cannot be evaluated.
	ErrGradeCodeMissing = errors.New("grade code is missing")

	// The current grade row is missing its interval, so the
	// interval rule cannot be evaluated.
	ErrGradeIntervalMissing = errors.New("grade interval is missing")

	// The target grade is not exactly one step above the
	// person's current grade.
	ErrGradeNotNext = errors.New("target grade must be exactly one step above the current grade")

	// The person has not served the interval of the current
	// grade yet.
	ErrGradeIntervalNotServed = errors.New("the grade interval of the current grade has not been served yet")
)

type Service struct {
	repo *Repository

	// The persons and pathGrades modules own their tables:
	// the candidate's person and target path grade are loaded
	// through their repositories, including the grade rules,
	// instead of this module re-implementing those lookups.
	personRepo    *persons.Repository
	pathGradeRepo *pathGrades.Repository
}

func NewService(repo *Repository, personRepo *persons.Repository, pathGradeRepo *pathGrades.Repository) *Service {
	return &Service{
		repo:          repo,
		personRepo:    personRepo,
		pathGradeRepo: pathGradeRepo,
	}
}

// loadPersonForValidation loads the candidate's person with the
// current grade preloaded, through the persons module repository
// scoped to the same transaction and client.
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

// loadPathGradeForValidation loads the target path grade with
// its grade preloaded, through the pathGrades module repository
// scoped to the same transaction and client.
func (s *Service) loadPathGradeForValidation(tx *gorm.DB, clientID int, pathGradeID int) (*models.PathGrade, error) {

	pathGrade, err := s.pathGradeRepo.
		WithDB(tx).
		ForClient(clientID).
		GetByID(pathGradeID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPathGradeNotFound
		}
		return nil, err
	}

	return pathGrade, nil
}

func (s *Service) GetAll(clientID int) ([]models.PathGradeCandidate, error) {

	repo := s.repo.ForClient(clientID)

	return repo.GetAll()
}

func (s *Service) GetByID(clientID int, id int) (*models.PathGradeCandidate, error) {

	repo := s.repo.ForClient(clientID)

	return repo.GetByID(id)
}

// Create inserts a new candidate for the authenticated client.
// The next-grade and interval rules of the person's current
// grade are checked inside the transaction, before the INSERT,
// so an invalid candidate aborts the whole transaction.
func (s *Service) Create(clientID int, req *models.PathGradeCandidate) (*models.PathGradeCandidate, error) {

	var candidate models.PathGradeCandidate

	err := dbutil.WithTransaction(s.repo.DB(), func(tx *gorm.DB) error {

		repo := s.repo.
			WithDB(tx).
			ForClient(clientID)

		// Load the person and the target path grade first:
		// both lookups double as existence checks, and their
		// grades feed the rules below.
		person, err := s.loadPersonForValidation(tx, clientID, req.PersonID)

		if err != nil {
			return err
		}

		pathGrade, err := s.loadPathGradeForValidation(tx, clientID, req.PathGradeID)

		if err != nil {
			return err
		}

		// Rule: the target grade must be exactly one step
		// above the person's current grade.
		if err := validateNextGrade(person, pathGrade); err != nil {
			return err
		}

		// Rule: the person must have served the interval of
		// the current grade, measured from the date of the
		// current grade to today.
		if err := validateGradeInterval(person); err != nil {
			return err
		}

		candidate = models.PathGradeCandidate{
			PersonID:          req.PersonID,
			PathGradeID:       req.PathGradeID,
			TermID:            req.TermID,
			FinalEvaluationID: req.FinalEvaluationID,
			ClientID:          clientID,
		}

		return repo.Create(&candidate)
	})

	if err != nil {
		return nil, err
	}

	return &candidate, nil
}

// Update is a full update of a candidate. The same rules of
// creation apply to the new values, and RowsAffected catches
// the missing-id case: a candidate of another client is
// reported exactly like a missing one, so ids cannot be
// probed across tenants.
func (s *Service) Update(clientID int, id int, req *models.PathGradeCandidate) (*models.PathGradeCandidate, error) {

	err := dbutil.WithTransaction(s.repo.DB(), func(tx *gorm.DB) error {

		repo := s.repo.
			WithDB(tx).
			ForClient(clientID)

		// The same rules of creation apply to the new
		// person and path grade values.
		person, err := s.loadPersonForValidation(tx, clientID, req.PersonID)

		if err != nil {
			return err
		}

		pathGrade, err := s.loadPathGradeForValidation(tx, clientID, req.PathGradeID)

		if err != nil {
			return err
		}

		if err := validateNextGrade(person, pathGrade); err != nil {
			return err
		}

		if err := validateGradeInterval(person); err != nil {
			return err
		}

		candidate := models.PathGradeCandidate{
			PersonID:          req.PersonID,
			PathGradeID:       req.PathGradeID,
			TermID:            req.TermID,
			FinalEvaluationID: req.FinalEvaluationID,
		}

		rows, err := repo.Update(id, &candidate)

		if err != nil {
			return err
		}

		if rows == 0 {
			return ErrPathGradeCandidateNotFound
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.repo.ForClient(clientID).GetByID(id)
}

// Delete removes a candidate by id. The client-scoped
// repository makes the id of another client's candidate behave
// exactly like a missing one, and RowsAffected catches both
// cases.
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
				return ErrPathGradeCandidateNotFound
			}

			return nil
		},
	)
}

// truncateToDate strips the time component of a time.Time so
// plain `date` columns can be compared day by day, exactly
// like the course session validations do.
func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// yearsSince counts whole years elapsed since the given date,
// truncated to whole days. A partial year does not count: 7
// years and 11 months in an 8-year grade is still 7.
func yearsSince(from time.Time, now time.Time) int {

	from = truncateToDate(from)
	now = truncateToDate(now)

	years := now.Year() - from.Year()

	// The birthday-style correction: subtract one if the
	// anniversary of `from` has not been reached yet this
	// year.
	anniversary := time.Date(
		from.Year()+years,
		from.Month(),
		from.Day(),
		0, 0, 0, 0, time.UTC,
	)

	if anniversary.After(now) {
		years--
	}

	return years
}

// validateNextGrade enforces the next-grade rule: the target
// grade of the path must be exactly one step above the person's
// current grade (3rd to 2nd, 2nd to 1st), comparing the numeric
// code of each grade. Grades with no code cannot be compared.
func validateNextGrade(person *models.Person, pathGrade *models.PathGrade) error {

	// The person must carry a current grade: candidates
	// without one cannot be placed on a grade ladder.
	if person.CurrentGradeID == nil || person.CurrentGrade == nil {
		return ErrPersonHasNoCurrentGrade
	}

	currentCode := person.CurrentGrade.Code
	targetCode := pathGrade.Grade.Code

	if currentCode == nil || targetCode == nil {
		return ErrGradeCodeMissing
	}

	if *targetCode != *currentCode-1 {
		return ErrGradeNotNext
	}

	return nil
}

// validateGradeInterval enforces the interval rule: the person
// must have served at least the interval of the current grade,
// measured from the date of the current grade to today. The
// interval belongs to the grade being left: it states how long
// one sits in it before moving to the next.
func validateGradeInterval(person *models.Person) error {

	// Already checked by validateNextGrade, but the
	// functions must stay independently safe.
	if person.CurrentGradeID == nil || person.CurrentGrade == nil {
		return ErrPersonHasNoCurrentGrade
	}

	gradeDate := person.DateOfCurrentGrade

	if gradeDate == nil {
		return ErrPersonHasNoGradeDate
	}

	interval := person.CurrentGrade.GradeInterval

	if interval == nil {
		return ErrGradeIntervalMissing
	}

	served := yearsSince(*gradeDate, time.Now())

	if served < *interval {
		return ErrGradeIntervalNotServed
	}

	return nil
}
