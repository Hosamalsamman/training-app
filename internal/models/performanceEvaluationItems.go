package models

type PerformanceEvaluationItem struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	Name string `gorm:"column:name;size:3000;not null" json:"name"`

	MinimumScore float64 `gorm:"column:minimum_score;type:numeric(5,2);not null" json:"minimum_score"`
	MaximumScore float64 `gorm:"column:maximum_score;type:numeric(5,2);not null" json:"maximum_score"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`
}

func (PerformanceEvaluationItem) TableName() string {
	return `"performanceEvaluation_items"`
}