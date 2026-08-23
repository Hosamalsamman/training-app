package models

type LearningPath struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Code *int   `gorm:"column:code" json:"code"`
	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	QualificationTypeID int `gorm:"column:qualification_type_id;not null" json:"qualification_type_id"`
	QualificationType QualificationType `gorm:"foreignKey:QualificationTypeID" json:"qualification_type"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`
}

func (LearningPath) TableName() string {
	return "learning_paths"
}