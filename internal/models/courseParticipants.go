package models

type CourseParticipant struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	CourseID int    `gorm:"column:course_id;not null" json:"course_id"`
	Course   Course `gorm:"foreignKey:CourseID;references:ID" json:"course"`

	PersonID int    `gorm:"column:person_id;not null" json:"person_id"`
	Person   Person `gorm:"foreignKey:PersonID;references:ID" json:"person"`

	FinalEvaluationID *int `gorm:"column:finalevaluation_id" json:"finalevaluation_id"`
	FinalEvaluation   *EvaluationType `gorm:"foreignKey:FinalEvaluationID;references:ID" json:"final_evaluation"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`

	FinalExams []CourseParticipantFinalExam `gorm:"foreignKey:CourseParticipantID;references:ID" json:"final_exams"`
}