package models

type PathGrade struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	LearningPathID int `gorm:"column:learning_path_id;not null" json:"learning_path_id"`
	LearningPath   LearningPath `gorm:"foreignKey:LearningPathID;references:ID" json:"learning_path"`

	GradeID int `gorm:"column:grade_id;not null" json:"grade_id"`
	Grade   Grade `gorm:"foreignKey:GradeID;references:ID" json:"grade"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`

	// Reverse relationship:
	// One PathGrade can have many PathGradeSubjects.
	PathGradeSubjects []PathGradeSubject `gorm:"foreignKey:PathGradeID;references:ID" json:"path_grade_subjects"`
}

func (PathGrade) TableName() string {
	return "path_grades"
}

