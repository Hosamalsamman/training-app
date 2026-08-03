package models

import "time"

type Course struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	// Either PathGradeSubjectTerm OR LearningSubject
	PathGradeSubjectTermID *int                  `gorm:"column:path_grade_subject_term_id" json:"path_grade_subject_term_id"`
	PathGradeSubjectTerm   *PathGradeSubjectTerm `gorm:"foreignKey:PathGradeSubjectTermID" json:"path_grade_subject_term,omitempty"`

	LearningSubjectID *int              `gorm:"column:learning_subject_id" json:"learning_subject_id"`
	LearningSubject   *LearningSubject  `gorm:"foreignKey:LearningSubjectID" json:"learning_subject,omitempty"`

	DurationInDays int       `gorm:"column:duration_in_days;not null" json:"duration_in_days"`
	StartingDate   time.Time `gorm:"column:starting_date;type:date;not null" json:"starting_date"`
	EndDate        time.Time `gorm:"column:end_date;type:date;not null" json:"end_date"`

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

	Cost float64 `gorm:"column:cost;type:numeric(9,3);not null" json:"cost"`

	IsPlanned  *bool `gorm:"column:is_palnned" json:"is_planned"`
	IsExecuted *bool `gorm:"column:is_executed" json:"is_executed"`

	PlannedID *int     `gorm:"column:planned_id" json:"planned_id"`
	Planned   *Course  `gorm:"foreignKey:PlannedID" json:"planned,omitempty"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`

	ExecutedCourses []Course `gorm:"foreignKey:PlannedID" json:"executed_courses,omitempty"`

	Sessions []CourseSession `gorm:"foreignKey:CourseID" json:"sessions,omitempty"`
}