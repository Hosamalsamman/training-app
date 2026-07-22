package models

type WorkGroup struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	Code string `gorm:"column:code;size:10" json:"code,omitempty"`

	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	ClientID int `gorm:"column:client_id;not null" json:"client_id"`

	JobTypeGroups []JobTypeGroup `gorm:"foreignKey:WorkGroupID" json:"job_type_groups,omitempty"`
}

func (WorkGroup) TableName() string {
	return "work_groups"
}