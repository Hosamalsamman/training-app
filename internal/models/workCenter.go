package models

type WorkCenter struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	GovernorateID int `gorm:"column:governorate_id;not null" json:"governorate_id"`
	Governorate Governorate `gorm:"foreignKey:GovernorateID" json:"governorate"`

	OrganizationID int `gorm:"column:organization_id;not null" json:"organization_id"`
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization"`

	ClientID int `gorm:"column:client_id;not null" json:"client_id"`

	WorkSites []WorkSite `gorm:"foreignKey:WorkCenterID" json:"work_sites,omitempty"`
}

func (WorkCenter) TableName() string {
	return "work_centers"
}