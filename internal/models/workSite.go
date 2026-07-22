package models

type WorkSite struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	Code string `gorm:"column:code;size:20" json:"code,omitempty"`

	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	WorkCenterID int `gorm:"column:work_center_id;not null" json:"work_center_id"`
	WorkCenter WorkCenter `gorm:"foreignKey:WorkCenterID" json:"work_center"`

	ClientID int `gorm:"column:client_id;not null" json:"client_id"`
}

func (WorkSite) TableName() string {
	return "work_sites"
}