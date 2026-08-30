package models

type PerformanceEvaluationParticipantCategoryItem struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	PerformanceEvaluationParticipantCategoryID int                                      `gorm:"column:performanceEvaluation_participant_category_id;not null" json:"performance_evaluation_participant_category_id"`
	PerformanceEvaluationParticipantCategory   PerformanceEvaluationParticipantCategory `gorm:"foreignKey:PerformanceEvaluationParticipantCategoryID;references:ID" json:"performance_evaluation_participant_category"`

	PerformanceEvaluationItemID int                       `gorm:"column:performanceEvaluation_item_id;not null" json:"performance_evaluation_item_id"`
	PerformanceEvaluationItem   PerformanceEvaluationItem `gorm:"foreignKey:PerformanceEvaluationItemID;references:ID" json:"performance_evaluation_item"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`

	CourseDetails []PerformanceEvaluationCourseDetail `gorm:"foreignKey:PerformanceEvaluationParticipantCategoryItemID;references:ID" json:"course_details,omitempty"`

	CourseParticipantDetails []PerformanceEvaluationCourseParticipantDetail `gorm:"foreignKey:PerformanceEvaluationParticipantCategoryItemID;references:ID" json:"course_participant_details,omitempty"`
}

func (PerformanceEvaluationParticipantCategoryItem) TableName() string {
	return `"performanceEvaluation_participant_category_items"`
}
