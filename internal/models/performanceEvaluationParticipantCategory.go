package models

type PerformanceEvaluationParticipantCategory struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	ParticipantTypeID int             `gorm:"column:participant_type_id;not null" json:"participant_type_id"`
	ParticipantType   ParticipantType `gorm:"foreignKey:ParticipantTypeID;references:ID" json:"participant_type"`

	PerformanceEvaluationCategoryID int                          `gorm:"column:performanceEvaluation_category_id;not null" json:"performance_evaluation_category_id"`
	PerformanceEvaluationCategory   PerformanceEvaluationCategory `gorm:"foreignKey:PerformanceEvaluationCategoryID;references:ID" json:"performance_evaluation_category"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`

	// One participant category can contain many evaluation items.
	Items []PerformanceEvaluationParticipantCategoryItem `gorm:"foreignKey:PerformanceEvaluationParticipantCategoryID;references:ID" json:"items,omitempty"`
}

func (PerformanceEvaluationParticipantCategory) TableName() string {
	return `"performanceEvaluation_participant_categoreis"`
}