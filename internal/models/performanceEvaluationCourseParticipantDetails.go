package models

import "time"

type PerformanceEvaluationCourseParticipantDetail struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	CourseParticipantID int               `gorm:"column:course_participant_id;not null" json:"course_participant_id"`
	CourseParticipant   CourseParticipant `gorm:"foreignKey:CourseParticipantID;references:ID" json:"course_participant"`

	Score float64 `gorm:"column:score;not null" json:"score"`

	EnteredBy       int    `gorm:"column:entered_by;not null" json:"entered_by"`
	EnteredByPerson Person `gorm:"foreignKey:EnteredBy;references:ID" json:"entered_by_person"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`

	EnteredAt time.Time `gorm:"column:entered_at;not null" json:"entered_at"`

	PerformanceEvaluationParticipantCategoryItemID int `gorm:"column:performanceEvaluation_participant_category_item_id;not null" json:"performance_evaluation_participant_category_item_id"`

	PerformanceEvaluationParticipantCategoryItem PerformanceEvaluationParticipantCategoryItem `gorm:"foreignKey:PerformanceEvaluationParticipantCategoryItemID;references:ID" json:"performance_evaluation_participant_category_item"`
}

func (PerformanceEvaluationCourseParticipantDetail) TableName() string {
	return `"performanceEvaluation_courseParticipant_details"`
}
