package models

type Documentation struct {
	ID       int    `gorm:"primaryKey;column:id" json:"id"`
	Name     string `gorm:"column:name;size:2000;not null" json:"name"`
	FilePath string `gorm:"column:file_path;size:4000;not null" json:"file_path"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`

	DocumentationTypeID *int               `gorm:"column:documentation_type_id" json:"documentation_type_id"`
	DocumentationType   *DocumentationType `gorm:"foreignKey:DocumentationTypeID" json:"documentation_type"`

	LearningSubjectID *int              `gorm:"column:learning_subject_id" json:"learning_subject_id"`
	LearningSubject   *LearningSubject  `gorm:"foreignKey:LearningSubjectID" json:"learning_subject"`

	PathGradeSubjectID *int              `gorm:"column:path_grade_subject_id" json:"path_grade_subject_id"`
	PathGradeSubject   *PathGradeSubject `gorm:"foreignKey:PathGradeSubjectID" json:"path_grade_subject"`

	PathGradeSubjectTermID *int                  `gorm:"column:path_grade_subject_term_id" json:"path_grade_subject_term_id"`
	PathGradeSubjectTerm   *PathGradeSubjectTerm `gorm:"foreignKey:PathGradeSubjectTermID" json:"path_grade_subject_term"`

	CourseID *int    `gorm:"column:course_id" json:"course_id"`
	Course   *Course `gorm:"foreignKey:CourseID" json:"course"`

	CourseSessionID *int           `gorm:"column:course_session_id" json:"course_session_id"`
	CourseSession   *CourseSession `gorm:"foreignKey:CourseSessionID" json:"course_session"`

	CourseParticipantFinalExams []CourseParticipantFinalExam `gorm:"foreignKey:DocumentationID;references:ID" json:"course_participant_final_exams"`
}
