package models

type CourseSessionParticipant struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	CourseSessionID int           `gorm:"column:course_session_id;not null" json:"course_session_id"`
	CourseSession   CourseSession `gorm:"foreignKey:CourseSessionID;references:ID" json:"course_session"`

	CourseParticipantID int               `gorm:"column:course_participant_id;not null" json:"course_participant_id"`
	CourseParticipant   CourseParticipant `gorm:"foreignKey:CourseParticipantID;references:ID" json:"course_participant"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`
}

// CourseSessionParticipantRequest is the payload for both
// POST /new-course-session-participant and PUT
// /course-session-participant/:id. The update is a full
// update, so its payload is identical to the create payload
// and both share this one struct. The row references the
// enrollment (course participant) instead of a raw person,
// so attendance can only be recorded for a person enrolled
// in a course. The service validates that the course session
// and the course participant belong to the same course.
type CourseSessionParticipantRequest struct {
	CourseSessionID     int `json:"course_session_id" binding:"required"`
	CourseParticipantID int `json:"course_participant_id" binding:"required"`
}