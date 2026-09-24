package models

type CourseParticipant struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	CourseID int    `gorm:"column:course_id;not null" json:"course_id"`
	Course   Course `gorm:"foreignKey:CourseID;references:ID" json:"course"`

	PersonID int    `gorm:"column:person_id;not null" json:"person_id"`
	Person   Person `gorm:"foreignKey:PersonID;references:ID" json:"person"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`

	FinalExams                   []CourseParticipantFinalExam                   `gorm:"foreignKey:CourseParticipantID;references:ID" json:"final_exams"`
	PerformanceEvaluationDetails []PerformanceEvaluationCourseParticipantDetail `gorm:"foreignKey:CourseParticipantID;references:ID" json:"performance_evaluation_details,omitempty"`
}

// CourseParticipantRequest is the payload for both
// POST /new-course-participant and PUT /course-participant/:id.
// The update is a full update, so its payload is identical to
// the create payload and both share this one struct. The
// service validates the path-grade rule: when the course is a
// path-grade subject (or subject term) course, the person must
// be a path grade candidate of the same path grade (and the
// same term for subject-term courses), with a final evaluation
// that is NULL or not a passing one
// (evaluation_types.should_repeat = 0), so an already passing
// person cannot take the course again.
type CourseParticipantRequest struct {
	CourseID int `json:"course_id" binding:"required"`
	PersonID int `json:"person_id" binding:"required"`
}

