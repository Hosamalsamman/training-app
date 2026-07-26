package models

type LearningSubject struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`

	TrainerSubjects []TrainerSubject `gorm:"foreignKey:LearningSubjectID" json:"trainer_subjects,omitempty"`
}

func (LearningSubject) TableName() string {
	return "learning_subjects"
}