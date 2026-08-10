package models

type PathGradeSubjectTerm struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	PathGradeSubjectID int `gorm:"column:path_grade_subject_id;not null" json:"path_grade_subjects_id"`
	PathGradeSubject PathGradeSubject `gorm:"foreignKey:PathGradeSubjectID;references:ID" json:"path_grade_subject"`

	TermID int          `gorm:"column:term_id;not null" json:"term_id"`
	Term   LearningTerm `gorm:"foreignKey:TermID" json:"term"`

	Code *int `gorm:"column:code" json:"code"`

	Notes *string `gorm:"column:notes;size:2000" json:"notes"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`

	Documentations []Documentation `gorm:"foreignKey:PathGradeSubjectTermID" json:"documentations"`
}

func (PathGradeSubjectTerm) TableName() string {
	return "path_grade_subject_terms"
}
