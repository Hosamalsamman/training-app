package models

import "time"

type CourseSession struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	SessionDate time.Time `gorm:"column:session_date;type:date;not null" json:"session_date"`
	StartTime   time.Time `gorm:"column:start_time;type:time;not null" json:"start_time"`
	EndTime     time.Time `gorm:"column:end_time;type:time;not null" json:"end_time"`

	CourseID int    `gorm:"column:course_id;not null" json:"course_id"`
	Course   Course `gorm:"foreignKey:CourseID" json:"course"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`

	Documentations []Documentation `gorm:"foreignKey:CourseSessionID" json:"documentations"`

	Participants []CourseSessionParticipant `gorm:"foreignKey:CourseSessionID;references:ID" json:"participants"`
}

// CourseSessionRequest is the payload for both
// POST /new-course-session and PUT /course-session/:id. The
// update is a full update, so its payload is identical to the
// create payload and both share this one struct. Every field
// is required so a missing field can never wipe existing data
// silently. The service validates that the session falls
// within the starting and ending dates of the course it
// belongs to, that start_time is before end_time, and that
// sessions are added to executed courses only.
type CourseSessionRequest struct {
	CourseID    int       `json:"course_id" binding:"required"`
	SessionDate time.Time `json:"session_date" binding:"required"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
}