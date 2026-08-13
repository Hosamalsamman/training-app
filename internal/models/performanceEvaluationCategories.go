package models

type PerformanceEvaluationCategory struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	Name string `gorm:"column:name;size:3000;not null" json:"name"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`

}

func (PerformanceEvaluationCategory) TableName() string {
	return `"performanceEvaluation_categories"`
}