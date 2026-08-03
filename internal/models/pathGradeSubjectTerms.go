package models

type PathGradeSubjectTerm struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	PathGradeSubjectsID int              `gorm:"column:path_grade_subjects_id;not null" json:"path_grade_subjects_id"`
	PathGradeSubject    PathGradeSubject `gorm:"foreignKey:PathGradeSubjectsID" json:"path_grade_subject"`

	TermID int          `gorm:"column:term_id;not null" json:"term_id"`
	Term   LearningTerm `gorm:"foreignKey:TermID" json:"term"`

	Code *int `gorm:"column:code" json:"code"`

	Notes *string `gorm:"column:notes;size:2000" json:"notes"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`
}

func (PathGradeSubjectTerm) TableName() string {
	return "path_grade_subject_term"
}