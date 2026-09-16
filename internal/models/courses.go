package models

import (
	"time"
	"github.com/shopspring/decimal"
)

type Course struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	// Either PathGradeSubjectTerm OR LearningSubject OR PathGradeSubject is not null
	PathGradeSubjectID *int               `gorm:"column:path_grade_subject_id" json:"path_grade_subject_id"`
	PathGradeSubject   *PathGradeSubject  `gorm:"foreignKey:PathGradeSubjectID" json:"path_grade_subject,omitempty"`
	
	PathGradeSubjectTermID *int                  `gorm:"column:path_grade_subject_term_id" json:"path_grade_subject_term_id"`
	PathGradeSubjectTerm   *PathGradeSubjectTerm `gorm:"foreignKey:PathGradeSubjectTermID" json:"path_grade_subject_term,omitempty"`

	LearningSubjectID *int              `gorm:"column:learning_subject_id" json:"learning_subject_id"`
	LearningSubject   *LearningSubject  `gorm:"foreignKey:LearningSubjectID" json:"learning_subject,omitempty"`

	DurationInDays int       `gorm:"column:duration_in_days;not null" json:"duration_in_days"`
	StartingDate   time.Time `gorm:"column:starting_date;type:date" json:"starting_date"`
	EndDate        time.Time `gorm:"column:end_date;type:date" json:"end_date"`

	NumberOfInternalParticipants *int `gorm:"column:number_of_internal_participants" json:"number_of_internal_participants"`
	NumberOfExternalParticipants *int `gorm:"column:number_of_external_participants" json:"number_of_external_participants"`

	RoomID int          `gorm:"column:room_id;not null" json:"room_id"`
	Room   TrainingRoom `gorm:"foreignKey:RoomID" json:"room"`

	FundingOrganizationID *int          `gorm:"column:funding_organization" json:"funding_organization_id"`
	FundingOrganization   *Organization `gorm:"foreignKey:FundingOrganizationID" json:"funding_organization,omitempty"`

	TrainerID int    `gorm:"column:trainer_id;not null" json:"trainer_id"`
	Trainer   Person `gorm:"foreignKey:TrainerID" json:"trainer"`

	BackupTrainerID *int    `gorm:"column:backup_trainer_id" json:"backup_trainer_id"`
	BackupTrainer   *Person `gorm:"foreignKey:BackupTrainerID" json:"backup_trainer,omitempty"`

	CoordinatorID *int    `gorm:"column:coordinator_id" json:"coordinator_id"`
	Coordinator   *Person `gorm:"foreignKey:CoordinatorID;references:ID" json:"coordinator,omitempty"`

	Cost decimal.Decimal `gorm:"column:cost;type:numeric(9,3);not null" json:"cost"`

	IsPlanned  *bool `gorm:"column:is_planned" json:"is_planned"`
	IsExecuted *bool `gorm:"column:is_executed" json:"is_executed"`

	PlannedID *int     `gorm:"column:planned_id" json:"planned_id"`
	Planned   *Course  `gorm:"foreignKey:PlannedID" json:"planned,omitempty"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`

	ExecutedCourses []Course `gorm:"foreignKey:PlannedID" json:"executed_courses,omitempty"`

	Sessions []CourseSession `gorm:"foreignKey:CourseID" json:"sessions,omitempty"`

	Documentations []Documentation `gorm:"foreignKey:CourseID" json:"documentations"`

	Participants []CourseParticipant `gorm:"foreignKey:CourseID;references:ID" json:"participants"`
}

// CreateCourseRequest is the payload for POST /courses.
// The service validates the two planning rules:
//
//  1. exactly one of the three subject identifiers is set,
//  2. is_planned / is_executed / planned_id form one of the
//     three valid course states.
type CreateCourseRequest struct {
	Name string `json:"name" binding:"required"`

	// Exactly one of these three identifiers must be set.
	PathGradeSubjectID     *int `json:"path_grade_subject_id"`
	PathGradeSubjectTermID *int `json:"path_grade_subject_term_id"`
	LearningSubjectID      *int `json:"learning_subject_id"`

	DurationInDays int       `json:"duration_in_days" binding:"required,min=1"`
	StartingDate   time.Time `json:"starting_date" binding:"required"`
	EndDate        time.Time `json:"end_date"`

	NumberOfInternalParticipants *int `json:"number_of_internal_participants"`
	NumberOfExternalParticipants *int `json:"number_of_external_participants"`

	RoomID    int `json:"room_id" binding:"required"`
	TrainerID int `json:"trainer_id" binding:"required"`

	FundingOrganizationID *int `json:"funding_organization_id"`
	BackupTrainerID       *int `json:"backup_trainer_id"`
	CoordinatorID         *int `json:"coordinator_id"`

	Cost decimal.Decimal `json:"cost"`

	// Lifecycle flags. Planned-not-executed,
	// executed-not-planned or executed-from-planned.
	IsPlanned  bool `json:"is_planned"`
	IsExecuted bool `json:"is_executed"`

	// Required only when the course is executed from a
	// planned course (is_planned && is_executed). The new row
	// carries the actual data, which may differ from the
	// planned course it was executed from.
	PlannedID *int `json:"planned_id"`
}

// CourseListFilters carries the optional subject query
// parameters of GET /get-planned-courses and
// GET /get-executed-courses. A nil field means the parameter
// was not sent, so the repository adds no WHERE condition
// for it.
type CourseListFilters struct {
	PathGradeSubjectID     *int
	PathGradeSubjectTermID *int
	LearningSubjectID      *int
}