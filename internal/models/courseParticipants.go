package models

type CourseParticipant struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	CourseID int    `gorm:"column:course_id;not null" json:"course_id"`
	Course   Course `gorm:"foreignKey:CourseID;references:ID" json:"course"`

	PersonID int    `gorm:"column:person_id;not null" json:"person_id"`
	Person   Person `gorm:"foreignKey:PersonID;references:ID" json:"person"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`

	FinalExams                   []CourseParticipantFinalExam                   `gorm:"foreignKey:CourseParticipantID;references:ID" json:"final_exams"`
	PerformanceEvaluationDetails []PerformanceEvaluationCourseParticipantDetail `gorm:"foreignKey:CourseParticipantID;references:ID" json:"performance_evaluation_details,omitempty"`
}
