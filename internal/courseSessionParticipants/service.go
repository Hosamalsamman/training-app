package courseSessionParticipants

import (
	"errors"
	"training-app/internal/courseParticipants"
	"training-app/internal/courseSessions"
	"training-app/internal/dbutil"
	"training-app/internal/models"

	"gorm.io/gorm"
)

// Sentinel errors for the attendance rules. The handler maps
// them to proper HTTP status codes.
var (
	// The referenced course session does not exist for this
	// client.
	ErrSessionNotFound = errors.New("course session not found for this client")

	// The referenced course participant (the enrollment) does
	// not exist for this client.
	ErrParticipantNotFound = errors.New("course participant not found for this client")

	// The attendance row id does not exist for this client.
	ErrSessionParticipantNotFound = errors.New("course session participant not found for this client")

	// The course participant is enrolled in another course
	// than the session: attendance can only be recorded for a
	// participant of the session's course.
	ErrParticipantNotInCourse = errors.New("the course participant is not enrolled in the course of this session")
)

type Service struct {
	repo *Repository

	// The courseSessions and courseParticipants modules own
	// their tables: the session and the enrollment an
	// attendance row references are loaded through their
	// repositories instead of this module re-implementing
	// those lookups.
	sessionRepo     *courseSessions.Repository
	participantRepo *courseParticipants.Repository
}

func NewService(repo *Repository, sessionRepo *courseSessions.Repository, participantRepo *courseParticipants.Repository) *Service {
	return &Service{
		repo:            repo,
		sessionRepo:     sessionRepo,
		participantRepo: participantRepo,
	}
}

func (s *Service) GetAll(clientID int) ([]models.CourseSessionParticipant, error) {

	repo := s.repo.ForClient(clientID)

	return repo.GetAll()
}

func (s *Service) GetByID(clientID int, id int) (*models.CourseSessionParticipant, error) {

	repo := s.repo.ForClient(clientID)

	return repo.GetByID(id)
}

// validateSameCourse enforces the attendance rule: the
// course session and the course participant (the enrollment)
// must belong to the same course, so attendance can only be
// recorded for a participant of the session's course. It is a
// pure function so it can be unit tested without a database.
func validateSameCourse(session *models.CourseSession, participant *models.CourseParticipant) error {

	if session.CourseID != participant.CourseID {
		return ErrParticipantNotInCourse
	}

	return nil
}

// Create validates the request, checks that the session and
// the enrollment belong to the same course, and inserts the
// attendance row, returning it with all relations preloaded.
func (s *Service) Create(clientID int, req *models.CourseSessionParticipantRequest) (*models.CourseSessionParticipant, error) {

	var participant models.CourseSessionParticipant

	err := dbutil.WithTransaction(
		s.repo.DB(),
		func(tx *gorm.DB) error {

			// Use the transaction DB and restrict this
			// repository to the authenticated client's data.
			repo := s.repo.
				WithDB(tx).
				ForClient(clientID)

			// The session must exist and belong to this
			// client. It is loaded through the
			// courseSessions module's repository scoped to
			// the same transaction, and the client filter
			// is the backend safety net: the foreign key
			// alone would silently accept another client's
			// session.
			session, err := s.sessionRepo.
				WithDB(tx).
				ForClient(clientID).
				GetByID(req.CourseSessionID)

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrSessionNotFound
				}
				return err
			}

			// The enrollment must exist and belong to this
			// client. Referencing the enrollment instead
			// of a raw person is what makes attendance
			// impossible for a person who was never
			// enrolled in a course.
			enrollment, err := s.participantRepo.
				WithDB(tx).
				ForClient(clientID).
				GetByID(req.CourseParticipantID)

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrParticipantNotFound
				}
				return err
			}

			// Rule: the session and the enrollment must
			// belong to the same course. The check runs
			// before the INSERT so an invalid pair aborts
			// the whole transaction.
			if err := validateSameCourse(session, enrollment); err != nil {
				return err
			}

			participant = models.CourseSessionParticipant{
				CourseSessionID:     req.CourseSessionID,
				CourseParticipantID: req.CourseParticipantID,
				ClientID:            clientID,
			}

			return repo.Create(&participant)
		},
	)

	if err != nil {
		return nil, err
	}

	return s.repo.ForClient(clientID).GetByID(participant.ID)
}

// Update re-applies the same validation of creation (both
// rows exist for this client and belong to the same course)
// to the new values and performs a full update, returning the
// row with all relations preloaded.
func (s *Service) Update(clientID int, id int, req *models.CourseSessionParticipantRequest) (*models.CourseSessionParticipant, error) {

	err := dbutil.WithTransaction(
		s.repo.DB(),
		func(tx *gorm.DB) error {

			repo := s.repo.
				WithDB(tx).
				ForClient(clientID)

			// The new session must exist and belong to
			// this client, so a row can never be moved to
			// another client's session.
			session, err := s.sessionRepo.
				WithDB(tx).
				ForClient(clientID).
				GetByID(req.CourseSessionID)

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrSessionNotFound
				}
				return err
			}

			// The new enrollment must exist and belong to
			// this client, so a row can never be moved to
			// another client's enrollment.
			enrollment, err := s.participantRepo.
				WithDB(tx).
				ForClient(clientID).
				GetByID(req.CourseParticipantID)

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrParticipantNotFound
				}
				return err
			}

			if err := validateSameCourse(session, enrollment); err != nil {
				return err
			}

			participant := models.CourseSessionParticipant{
				CourseSessionID:     req.CourseSessionID,
				CourseParticipantID: req.CourseParticipantID,
				ClientID:            clientID,
			}

			rows, err := repo.Update(id, &participant)

			if err != nil {
				return err
			}

			if rows == 0 {
				return ErrSessionParticipantNotFound
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return s.repo.ForClient(clientID).GetByID(id)
}

// Delete removes an attendance row by id. The client-scoped
// repository makes the id of another client's row behave
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
				return ErrSessionParticipantNotFound
			}

			return nil
		},
	)
}