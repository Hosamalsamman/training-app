package models

type Governorate struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	CountryID int `gorm:"column:country_id;not null" json:"country_id"`
	Country Country `gorm:"foreignKey:CountryID" json:"country"`

	ClientID int `gorm:"column:client_id;not null" json:"client_id"`

	Organizations []Organization `gorm:"foreignKey:GovernorateID" json:"organizations,omitempty"`

	WorkCenters []WorkCenter `gorm:"foreignKey:GovernorateID" json:"work_centers,omitempty"`
}

func (Governorate) TableName() string {
	return "governorates"
}
