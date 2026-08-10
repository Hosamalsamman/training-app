package models

import "time"

type CourseParticipantFinalExam struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	ExamDate time.Time `gorm:"column:exam_date;not null" json:"exam_date"`

	OrganizedBy int          `gorm:"column:organized_by;not null" json:"organized_by"`
	Organizer   Organization `gorm:"foreignKey:OrganizedBy;references:ID" json:"organizer"`

	CourseParticipantID int               `gorm:"column:course_participant_id;not null" json:"course_participant_id"`
	CourseParticipant   CourseParticipant `gorm:"foreignKey:CourseParticipantID;references:ID" json:"course_participant"`

	ExamScore float64 `gorm:"column:exam_score;type:numeric(5,2);not null" json:"exam_score"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`
}

func (CourseParticipantFinalExam) TableName() string {
	return `"courseParticipant_finalExams"`
}