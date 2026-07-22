package models

type JobTypeGroup struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	Code string `gorm:"column:code;size:10" json:"code,omitempty"`

	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	WorkGroupID int `gorm:"column:work_group_id;not null" json:"work_group_id"`
	WorkGroup WorkGroup `gorm:"foreignKey:WorkGroupID" json:"work_group"`

	ClientID int `gorm:"column:client_id;not null" json:"client_id"`
}

func (JobTypeGroup) TableName() string {
	return "job_type_groups"
}