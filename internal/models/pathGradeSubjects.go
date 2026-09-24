package models

type PathGradeSubject struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	Name string `gorm:"column:name;not null" json:"name"`

	Code *string `gorm:"column:code;size:10" json:"code"`

	PathGradeID int       `gorm:"column:path_grade_id;not null" json:"path_grade_id"`
	PathGrade   PathGrade `gorm:"foreignKey:PathGradeID;references:ID" json:"path_grade"`

	LearningSubjectID int             `gorm:"column:learning_subject_id;not null" json:"learning_subject_id"`
	LearningSubject   LearningSubject `gorm:"foreignKey:LearningSubjectID;references:ID" json:"learning_subject"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`

	PathGradeSubjectTerms []PathGradeSubjectTerm `gorm:"foreignKey:PathGradeSubjectID;references:ID" json:"path_grade_subject_terms"`

	Documentations []Documentation `gorm:"foreignKey:PathGradeSubjectID" json:"documentations"`
}

func (PathGradeSubject) TableName() string {
	return "path_grade_subjects"
}