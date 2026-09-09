package courses

import (
	"errors"
	"testing"

	"training-app/internal/models"
)

func TestValidateCreateRequest(t *testing.T) {

	// Subject identifiers used across the cases.
	pgsID := 1
	pgstID := 2
	lsID := 3

	// An existing pending planned course to reference.
	plannedID := 10

	tests := []struct {
		name    string
		req     models.CreateCourseRequest
		wantErr error
	}{
		{
			name: "planned not executed with path grade subject",
			req: models.CreateCourseRequest{
				Name:               "Course A",
				PathGradeSubjectID: &pgsID,
				IsPlanned:          true,
			},
		},
		{
			name: "executed not planned with learning subject",
			req: models.CreateCourseRequest{
				Name:              "Course B",
				LearningSubjectID: &lsID,
				IsExecuted:        true,
			},
		},
		{
			name: "executed from planned with subject term",
			req: models.CreateCourseRequest{
				Name:                   "Course C",
				PathGradeSubjectTermID: &pgstID,
				IsPlanned:              true,
				IsExecuted:             true,
				PlannedID:              &plannedID,
			},
		},
		{
			name: "no subject identifier set",
			req: models.CreateCourseRequest{
				Name:      "Course D",
				IsPlanned: true,
			},
			wantErr: ErrNoSubjectSet,
		},
		{
			name: "two subject identifiers set",
			req: models.CreateCourseRequest{
				Name:                   "Course E",
				PathGradeSubjectID:     &pgsID,
				PathGradeSubjectTermID: &pgstID,
				IsPlanned:              true,
			},
			wantErr: ErrMultipleSubjectsSet,
		},
		{
			name: "all three subject identifiers set",
			req: models.CreateCourseRequest{
				Name:                   "Course F",
				PathGradeSubjectID:     &pgsID,
				PathGradeSubjectTermID: &pgstID,
				LearningSubjectID:      &lsID,
				IsPlanned:              true,
			},
			wantErr: ErrMultipleSubjectsSet,
		},
		{
			name: "neither planned nor executed",
			req: models.CreateCourseRequest{
				Name:               "Course G",
				PathGradeSubjectID: &pgsID,
			},
			wantErr: ErrInvalidState,
		},
		{
			name: "planned not executed with planned_id",
			req: models.CreateCourseRequest{
				Name:               "Course H",
				PathGradeSubjectID: &pgsID,
				IsPlanned:          true,
				PlannedID:          &plannedID,
			},
			wantErr: ErrUnexpectedPlannedID,
		},
		{
			name: "executed not planned with planned_id",
			req: models.CreateCourseRequest{
				Name:               "Course I",
				PathGradeSubjectID: &pgsID,
				IsExecuted:         true,
				PlannedID:          &plannedID,
			},
			wantErr: ErrUnexpectedPlannedID,
		},
		{
			name: "executed from planned without planned_id",
			req: models.CreateCourseRequest{
				Name:               "Course J",
				PathGradeSubjectID: &pgsID,
				IsPlanned:          true,
				IsExecuted:         true,
			},
			wantErr: ErrPlannedIDRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := validateCreateRequest(&tt.req)

			// errors.Is(nil, nil) is true, so the valid
			// cases (wantErr nil) are covered too.
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("validateCreateRequest() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
