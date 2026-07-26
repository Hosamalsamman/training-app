package models

type TrainerSubject struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	PersonID int    `gorm:"column:person_id;not null" json:"person_id"`
	Person   Person `gorm:"foreignKey:PersonID" json:"person"`

	LearningSubjectID int             `gorm:"column:learning_subject_id;not null" json:"learning_subject_id"`
	LearningSubject   LearningSubject `gorm:"foreignKey:LearningSubjectID" json:"learning_subject"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`
}

func (TrainerSubject) TableName() string {
	return "trainer_subjects"
}