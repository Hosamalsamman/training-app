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
}