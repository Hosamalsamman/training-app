package courseSessions

import (
	"errors"
	"time"
	"training-app/internal/courses"
	"training-app/internal/dbutil"
	"training-app/internal/models"

	"gorm.io/gorm"
)

// Sentinel errors for the session rules. The handler maps
// them to proper HTTP status codes.
var (
	// end_time is not after start_time.
	ErrInvalidSessionTimes = errors.New("session start time must be before end time")

	// session_date falls outside the course date range.
	ErrSessionOutsideCourseDates = errors.New("session date must be within the course starting and ending dates")

	// The referenced course exists but is not executed:
	// sessions can be added to executed courses only.
	ErrCourseNotExecuted = errors.New("course sessions can be added to executed courses only")

	// The referenced course does not exist for this client.
	ErrCourseNotFound = errors.New("course not found for this client")

	// The session id does not exist for this client.
	ErrSessionNotFound = errors.New("course session not found for this client")
)

type Service struct {
	repo *Repository

	// The courses module owns its table: the course a
	// session belongs to is loaded through its repository,
	// including the executed-only rule, instead of this
	// module re-implementing a course lookup.
	courseRepo *courses.Repository
}

func NewService(repo *Repository, courseRepo *courses.Repository) *Service {
	return &Service{
		repo:       repo,
		courseRepo: courseRepo,
	}
}

func (s *Service) GetAll(clientID int) ([]models.CourseSession, error) {
	repo := s.repo.ForClient(clientID)
	return repo.GetAll()
}

func (s *Service) GetByID(clientID int, id int) (*models.CourseSession, error) {
	repo := s.repo.ForClient(clientID)
	return repo.GetByID(id)
}

// Create validates the request, checks the session against
// its course (executed only, session date within the course
// date range) and inserts it, returning it with all relations
// preloaded.
func (s *Service) Create(clientID int, req *models.CourseSessionRequest) (*models.CourseSession, error) {

	if err := validateSessionTimes(req); err != nil {
		return nil, err
	}

	var session models.CourseSession

	err := dbutil.WithTransaction(
		s.repo.DB(),
		func(tx *gorm.DB) error {

			// Use the transaction DB and restrict this
			// repository to the authenticated client's data.
			repo := s.repo.
				WithDB(tx).
				ForClient(clientID)

			// The course must exist and belong to this
			// client. It is loaded through the courses
			// module's repository scoped to the same
			// transaction, and the client filter is the
			// backend safety net: the course_id foreign
			// key alone would silently accept another
			// client's course.
			course, err := s.courseRepo.
				WithDB(tx).
				ForClient(clientID).
				GetByID(req.CourseID)

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrCourseNotFound
				}
				return err
			}

			// Rule: sessions can be added to executed
			// courses only. The check runs before the
			// INSERT so an invalid course aborts the
			// whole transaction.
			if err := validateCourseIsExecuted(course); err != nil {
				return err
			}

			// Rule: the session date must be within the
			// course starting and ending dates. The check
			// runs before the INSERT so an invalid date
			// aborts the whole transaction.
			if err := validateSessionWithinCourse(req.SessionDate, course); err != nil {
				return err
			}

			// Build the database model from the request.
			session = models.CourseSession{
				SessionDate: req.SessionDate,
				StartTime:   req.StartTime,
				EndTime:     req.EndTime,
				CourseID:    req.CourseID,

				// IMPORTANT:
				// client_id comes from JWT, not the request.
				ClientID: clientID,
			}

			// Repository only performs INSERT.
			// It does not commit.
			return repo.Create(&session)
		},
	)

	if err != nil {
		return nil, err
	}

	// Re-fetch with preloads so the response contains the
	// related data exactly like GET /course-session/:id does.
	return s.repo.ForClient(clientID).GetByID(session.ID)
}

// Update performs a full update of an existing session. The
// same rules of creation apply: the course referenced by
// course_id is loaded again, must be executed, and the new
// session_date is checked against the actual course dates.
func (s *Service) Update(clientID int, id int, req *models.CourseSessionRequest) (*models.CourseSession, error) {

	// The update payload carries the same fields as the
	// create payload, so the time-order rule is checked
	// through the same helper on the same struct.
	if err := validateSessionTimes(req); err != nil {
		return nil, err
	}

	err := dbutil.WithTransaction(
		s.repo.DB(),
		func(tx *gorm.DB) error {

			repo := s.repo.
				WithDB(tx).
				ForClient(clientID)

			// The session must exist and belong to this
			// client. The client-scoped WHERE makes other
			// clients' sessions invisible here.
			if _, err := repo.GetByID(id); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrSessionNotFound
				}
				return err
			}

			// The new course must exist, belong to this
			// client and be executed, so a session can
			// never be moved to another client's course
			// or to a course that was never executed.
			course, err := s.courseRepo.
				WithDB(tx).
				ForClient(clientID).
				GetByID(req.CourseID)

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrCourseNotFound
				}
				return err
			}

			// Rule: sessions can be added to executed
			// courses only.
			if err := validateCourseIsExecuted(course); err != nil {
				return err
			}

			if err := validateSessionWithinCourse(req.SessionDate, course); err != nil {
				return err
			}

			session := models.CourseSession{
				SessionDate: req.SessionDate,
				StartTime:   req.StartTime,
				EndTime:     req.EndTime,
				CourseID:    req.CourseID,
				ClientID:    clientID,
			}

			rows, err := repo.Update(id, &session)

			if err != nil {
				return err
			}

			if rows == 0 {
				return ErrSessionNotFound
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return s.repo.ForClient(clientID).GetByID(id)
}

// Delete removes a session by id. The client-scoped repository
// makes the id of another client's session behave exactly like
// a missing one, and RowsAffected catches both cases.
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
				return ErrSessionNotFound
			}

			return nil
		},
	)
}

// truncateToDate strips the time component of a time.Time so
// plain `date` columns can be compared day by day. Both the
// course dates and the session date are stored without a time
// zone, so the location of the value itself is preserved.
func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// validateSessionTimes enforces the time-order rule without
// touching the database: start_time must be before end_time.
func validateSessionTimes(req *models.CourseSessionRequest) error {

	if !req.StartTime.Before(req.EndTime) {
		return ErrInvalidSessionTimes
	}

	return nil
}

// validateCourseIsExecuted enforces the executed-only rule:
// sessions can be added to executed courses only. IsExecuted
// is a nullable boolean in the database, so a NULL counts as
// not executed.
func validateCourseIsExecuted(course *models.Course) error {

	if course.IsExecuted == nil || !*course.IsExecuted {
		return ErrCourseNotExecuted
	}

	return nil
}

// validateSessionWithinCourse enforces the date-range rule:
// session_date must fall within the starting and ending dates
// of the course it belongs to. Both sides are truncated to
// whole days first, since the columns are plain `date` types.
func validateSessionWithinCourse(sessionDate time.Time, course *models.Course) error {

	date := truncateToDate(sessionDate)
	start := truncateToDate(course.StartingDate)
	end := truncateToDate(course.EndDate)

	if date.Before(start) || date.After(end) {
		return ErrSessionOutsideCourseDates
	}

	return nil
}