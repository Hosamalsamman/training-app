package models

type Grade struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Code *int   `gorm:"column:code" json:"code"`
	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	GradeInterval *int `gorm:"column:grade_interval_in_years" json:"grade_interval"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`
}

func (Grade) TableName() string {
	return "grades"
}