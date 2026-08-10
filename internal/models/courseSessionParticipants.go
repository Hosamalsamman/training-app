package models

type CourseSessionParticipant struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	CourseSessionID int           `gorm:"column:course_session_id;not null" json:"course_session_id"`
	CourseSession   CourseSession `gorm:"foreignKey:CourseSessionID;references:ID" json:"course_session"`

	PersonID int    `gorm:"column:person_id;not null" json:"person_id"`
	Person   Person `gorm:"foreignKey:PersonID;references:ID" json:"person"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID;references:ID" json:"client"`
}