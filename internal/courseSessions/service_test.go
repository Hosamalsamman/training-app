package courseSessions

import (
	"errors"
	"testing"
	"time"
	"training-app/internal/models"
)

// A course running from the 1st to the 10th of September 2026,
// used as the date-range reference across the cases.
func testCourse() *models.Course {
	return &models.Course{
		StartingDate: time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2026, time.September, 10, 0, 0, 0, 0, time.UTC),
	}
}

func TestValidateSessionTimes(t *testing.T) {

	tests := []struct {
		name    string
		start   time.Time
		end     time.Time
		wantErr error
	}{
		{
			name:  "start before end",
			start: time.Date(0, time.January, 1, 9, 0, 0, 0, time.UTC),
			end:   time.Date(0, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:    "start after end",
			start:   time.Date(0, time.January, 1, 15, 0, 0, 0, time.UTC),
			end:     time.Date(0, time.January, 1, 12, 0, 0, 0, time.UTC),
			wantErr: ErrInvalidSessionTimes,
		},
		{
			name:    "start equals end",
			start:   time.Date(0, time.January, 1, 12, 0, 0, 0, time.UTC),
			end:     time.Date(0, time.January, 1, 12, 0, 0, 0, time.UTC),
			wantErr: ErrInvalidSessionTimes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := &models.CourseSessionRequest{
				CourseID:  1,
				StartTime: tt.start,
				EndTime:   tt.end,
			}

			err := validateSessionTimes(req)

			// errors.Is(nil, nil) is true, so the valid
			// cases (wantErr nil) are covered too.
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("validateSessionTimes() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSessionWithinCourse(t *testing.T) {

	course := testCourse()

	tests := []struct {
		name        string
		sessionDate time.Time
		wantErr     error
	}{
		{
			name:        "first day of the course",
			sessionDate: time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "last day of the course",
			sessionDate: time.Date(2026, time.September, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "middle of the course",
			sessionDate: time.Date(2026, time.September, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "day before the course starts",
			sessionDate: time.Date(2026, time.August, 31, 0, 0, 0, 0, time.UTC),
			wantErr:     ErrSessionOutsideCourseDates,
		},
		{
			name:        "day after the course ends",
			sessionDate: time.Date(2026, time.September, 11, 0, 0, 0, 0, time.UTC),
			wantErr:     ErrSessionOutsideCourseDates,
		},
		{
			// The session date carries a time component
			// because the frontend sent a full timestamp:
			// the day must still count as within range.
			name:        "same day with a time component",
			sessionDate: time.Date(2026, time.September, 10, 15, 30, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := validateSessionWithinCourse(tt.sessionDate, course)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("validateSessionWithinCourse() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCourseIsExecuted(t *testing.T) {

	// Pointers shared by the table cases; the helper never
	// mutates them.
	executed := true
	notExecuted := false

	tests := []struct {
		name    string
		course  *models.Course
		wantErr error
	}{
		{
			name:   "executed course",
			course: &models.Course{IsExecuted: &executed},
		},
		{
			name:    "planned course not executed yet",
			course:  &models.Course{IsExecuted: &notExecuted},
			wantErr: ErrCourseNotExecuted,
		},
		{
			// is_executed is nullable in the database:
			// a NULL row counts as not executed.
			name:    "null is_executed",
			course:  &models.Course{IsExecuted: nil},
			wantErr: ErrCourseNotExecuted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := validateCourseIsExecuted(tt.course)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("validateCourseIsExecuted() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}