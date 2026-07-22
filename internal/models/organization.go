package models

type Organization struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	Code string `gorm:"column:code;size:20" json:"code,omitempty"`

	Address string `gorm:"column:address;size:2000" json:"address,omitempty"`

	OrganizationTypeID int `gorm:"column:organization_type_id;not null" json:"organization_type_id"`
	OrganizationType OrganizationType `gorm:"foreignKey:OrganizationTypeID" json:"organization_type"`

	GovernorateID int `gorm:"column:governorate_id;not null" json:"governorate_id"`
	Governorate Governorate `gorm:"foreignKey:GovernorateID" json:"governorate"`

	ClientID int `gorm:"column:client_id;not null" json:"client_id"`

	WorkCenters []WorkCenter `gorm:"foreignKey:OrganizationID" json:"work_centers,omitempty"`
}

func (Organization) TableName() string {
	return "organizations"
}