package models

type PerformanceEvaluationParticipantCategoryItem struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	ParticipantTypeID int             `gorm:"column:participant_type_id;not null" json:"participant_type_id"`
	ParticipantType   ParticipantType `gorm:"foreignKey:ParticipantTypeID;references:ID" json:"participant_type"`

	PerformanceEvaluationCategoryID int                        `gorm:"column:performanceEvaluation_category_id;not null" json:"performance_evaluation_category_id"`
	PerformanceEvaluationCategory   PerformanceEvaluationCategory `gorm:"foreignKey:PerformanceEvaluationCategoryID;references:ID" json:"performance_evaluation_category"`

	PerformanceEvaluationItemID int                    `gorm:"column:performanceEvaluation_item_id;not null" json:"performance_evaluation_item_id"`
	PerformanceEvaluationItem   PerformanceEvaluationItem `gorm:"foreignKey:PerformanceEvaluationItemID;references:ID" json:"performance_evaluation_item"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`
}

func (PerformanceEvaluationParticipantCategoryItem) TableName() string {
	return `"performanceEvaluation_participant_category_items"`
}