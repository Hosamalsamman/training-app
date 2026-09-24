package models

type PathGradeCandidate struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	PersonID int    `gorm:"column:person_id;not null" json:"person_id"`
	Person   Person `gorm:"foreignKey:PersonID;references:ID" json:"person"`

	PathGradeID int       `gorm:"column:path_grade_id;not null" json:"path_grade_id"`
	PathGrade   PathGrade `gorm:"foreignKey:PathGradeID;references:ID" json:"path_grade"`

	TermID int          `gorm:"column:term_id;not null" json:"term_id"`
	Term   LearningTerm `gorm:"foreignKey:TermID;references:ID" json:"term"`

	FinalEvaluationID *int            `gorm:"column:final_evaluation" json:"final_evaluation_id"`
	FinalEvaluation   *EvaluationType `gorm:"foreignKey:FinalEvaluationID;references:ID" json:"final_evaluation"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`
}

func (PathGradeCandidate) TableName() string {
	return "path_grade_candidates"
}