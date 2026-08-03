package models

type PathGradeSubject struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	LearningPathID int          `gorm:"column:learning_path_id;not null" json:"learning_path_id"`
	LearningPath   LearningPath `gorm:"foreignKey:LearningPathID" json:"learning_path"`

	GradeID int   `gorm:"column:grade_id;not null" json:"grade_id"`
	Grade   Grade `gorm:"foreignKey:GradeID" json:"grade"`

	LearningSubjectID int             `gorm:"column:learning_subject_id;not null" json:"learning_subject_id"`
	LearningSubject   LearningSubject `gorm:"foreignKey:LearningSubjectID" json:"learning_subject"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`

	PathGradeSubjectTerms []PathGradeSubjectTerm `gorm:"foreignKey:PathGradeSubjectID" json:"path_grade_subject_terms,omitempty"`
}

func (PathGradeSubject) TableName() string {
	return "path_grade_subjects"
}