package courses

import (
	"errors"
	"training-app/internal/dbutil"
	"training-app/internal/models"

	"gorm.io/gorm"
)

// Sentinel errors for the course planning rules.
// The handler maps them to proper HTTP status codes.
var (
	// No subject identifier was set at all.
	ErrNoSubjectSet = errors.New("course must have exactly one of path_grade_subject_id, path_grade_subject_term_id or learning_subject_id")

	// More than one subject identifier was set.
	ErrMultipleSubjectsSet = errors.New("course must not have more than one of path_grade_subject_id, path_grade_subject_term_id or learning_subject_id")

	// The course is neither planned nor executed.
	ErrInvalidState = errors.New("course must be planned, executed or executed from a planned course")

	// planned_id was sent for a course that is not executed
	// from a planned course.
	ErrUnexpectedPlannedID = errors.New("planned_id is only allowed when the course is executed from a planned course")

	// planned_id is missing for a course executed from a
	// planned course.
	ErrPlannedIDRequired = errors.New("planned_id is required when the course is executed from a planned course")

	// The referenced planned course does not exist for this
	// client or was already executed (a planned course can
	// be executed only once).
	ErrPlannedCourseNotFound = errors.New("planned course not found for this client")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(clientID int) ([]models.Course, error) {
	repo := s.repo.ForClient(clientID)
	return repo.GetAll()
}

func (s *Service) GetByID(clientID int, id int) (*models.Course, error) {
	repo := s.repo.ForClient(clientID)
	return repo.GetByID(id)
}

// GetPlanned returns the client's pending planned courses
// (is_planned = true AND is_executed = false), optionally
// narrowed by the subject filters the frontend sent as
// query parameters.
func (s *Service) GetPlanned(clientID int, filters models.CourseListFilters) ([]models.Course, error) {
	repo := s.repo.ForClient(clientID)
	return repo.GetAllPlanned(filters)
}

// GetExecuted returns the client's executed courses
// (is_executed = true), optionally narrowed by the subject
// filters the frontend sent as query parameters. Sessions can
// be added to executed courses only, so these are the courses
// the frontend renders in the session course picker.
func (s *Service) GetExecuted(clientID int, filters models.CourseListFilters) ([]models.Course, error) {
	repo := s.repo.ForClient(clientID)
	return repo.GetAllExecuted(filters)
}

// validateCreateRequest enforces the two business rules of
// course planning without touching the database:
//
// Rule 1: exactly one of PathGradeSubjectID,
// PathGradeSubjectTermID or LearningSubjectID must be set.
//
// Rule 2: is_planned / is_executed must form one of the
// three valid states, and planned_id is only allowed when
// the course is executed from a planned course.
func validateCreateRequest(req *models.CreateCourseRequest) error {

	// Rule 1: count how many subject identifiers are set.
	set := 0

	if req.PathGradeSubjectID != nil {
		set++
	}

	if req.PathGradeSubjectTermID != nil {
		set++
	}

	if req.LearningSubjectID != nil {
		set++
	}

	switch {
	case set == 0:
		return ErrNoSubjectSet
	case set > 1:
		return ErrMultipleSubjectsSet
	}

	// Rule 2: lifecycle state and planned_id consistency.
	switch {

	case req.IsPlanned && req.IsExecuted:
		// Executed from a planned course: planned_id required.
		if req.PlannedID == nil {
			return ErrPlannedIDRequired
		}

	case req.IsPlanned || req.IsExecuted:
		// Planned-not-executed or executed-not-planned:
		// there is no planned course to reference.
		if req.PlannedID != nil {
			return ErrUnexpectedPlannedID
		}

	default:
		// Neither planned nor executed is meaningless.
		return ErrInvalidState
	}

	return nil
}

// Create validates the request, inserts the course and
// returns it with all relations preloaded.
func (s *Service) Create(clientID int, req *models.CreateCourseRequest) (*models.Course, error) {

	if err := validateCreateRequest(req); err != nil {
		return nil, err
	}

	var course models.Course

	err := dbutil.WithTransaction(
		s.repo.DB(),
		func(tx *gorm.DB) error {

			// Use the transaction DB and restrict this
			// repository to the authenticated client's data.
			repo := s.repo.
				WithDB(tx).
				ForClient(clientID)

			// Executed from a planned course: the referenced
			// course must exist, belong to this client and
			// not be executed yet. A planned course can be
			// executed only once, and "already executed" is
			// visible only through the planned_id column, so
			// GetPlannableCourse rejects a second execution
			// even if the frontend sent a stale id. The check
			// runs before the INSERT so an invalid planned
			// course aborts the whole transaction.
			if req.PlannedID != nil {
				if _, err := repo.GetPlannableCourse(*req.PlannedID); err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return ErrPlannedCourseNotFound
					}
					return err
				}
			}

			// Build the database model from the request.
			course = models.Course{
				Name:                         req.Name,
				PathGradeSubjectID:           req.PathGradeSubjectID,
				PathGradeSubjectTermID:       req.PathGradeSubjectTermID,
				LearningSubjectID:            req.LearningSubjectID,
				DurationInDays:               req.DurationInDays,
				StartingDate:                 req.StartingDate,
				EndDate:                      req.EndDate,
				NumberOfInternalParticipants: req.NumberOfInternalParticipants,
				NumberOfExternalParticipants: req.NumberOfExternalParticipants,
				RoomID:                       req.RoomID,
				FundingOrganizationID:        req.FundingOrganizationID,
				TrainerID:                    req.TrainerID,
				BackupTrainerID:              req.BackupTrainerID,
				CoordinatorID:                req.CoordinatorID,
				Cost:                         req.Cost,
				IsPlanned:                    &req.IsPlanned,
				IsExecuted:                   &req.IsExecuted,
				PlannedID:                    req.PlannedID,

				// IMPORTANT:
				// client_id comes from JWT, not the request.
				ClientID: clientID,
			}

			// Repository only performs INSERT.
			// It does not commit.
			return repo.Create(&course)
		},
	)

	if err != nil {
		return nil, err
	}

	// Re-fetch with preloads so the response contains the
	// related data exactly like GET /course/:id does.
	return s.repo.ForClient(clientID).GetByID(course.ID)
}
