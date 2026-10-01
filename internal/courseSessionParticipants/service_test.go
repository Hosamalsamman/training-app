package courseSessionParticipants

import (
	"errors"
	"testing"

	"training-app/internal/models"
)

func TestValidateSameCourse(t *testing.T) {

	tests := []struct {
		name        string
		session     *models.CourseSession
		participant *models.CourseParticipant
		wantErr     error
	}{
		{
			name:        "enrollment of the same course as the session",
			session:     &models.CourseSession{CourseID: 1},
			participant: &models.CourseParticipant{CourseID: 1},
		},
		{
			name:        "enrollment of another course",
			session:     &models.CourseSession{CourseID: 1},
			participant: &models.CourseParticipant{CourseID: 2},
			wantErr:     ErrParticipantNotInCourse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := validateSameCourse(tt.session, tt.participant)

			// errors.Is(nil, nil) is true, so the valid
			// cases (wantErr nil) are covered too.
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("validateSameCourse() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}