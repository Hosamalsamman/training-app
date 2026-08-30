package models

import (
	"time"
)

type PerformanceEvaluationCourseDetail struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	CourseID int    `gorm:"column:course_id;not null" json:"course_id"`
	Course   Course `gorm:"foreignKey:CourseID;references:ID" json:"course"`

	Score float64 `gorm:"column:score;not null" json:"score"`

	Notes *string `gorm:"column:notes;size:4000" json:"notes"`

	EnteredByID int    `gorm:"column:entered_by;not null" json:"entered_by"`
	EnteredBy   Person `gorm:"foreignKey:EnteredByID;references:ID" json:"entered_by_person"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`

	EnteredAt time.Time `gorm:"column:entered_at;type:date;not null" json:"entered_at"`

	PerformanceEvaluationParticipantCategoryItemID int                                          `gorm:"column:performanceEvaluation_participant_category_item_id;not null" json:"performance_evaluation_participant_category_item_id"`
	PerformanceEvaluationParticipantCategoryItem   PerformanceEvaluationParticipantCategoryItem `gorm:"foreignKey:PerformanceEvaluationParticipantCategoryItemID;references:ID" json:"performance_evaluation_participant_category_item"`
}

func (PerformanceEvaluationCourseDetail) TableName() string {
	return `"performanceEvaluation_course_details"`
}
