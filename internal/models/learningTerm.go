package models

type LearningTerm struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`

	PathGradeSubjectTerms []PathGradeSubjectTerm `gorm:"foreignKey:TermID" json:"path_grade_subject_terms,omitempty"`
}