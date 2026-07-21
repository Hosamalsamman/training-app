package models

type Governorate struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	CountryID int     `gorm:"column:country_id;not null" json:"country_id"`
	Country   Country `gorm:"foreignKey:CountryID" json:"country"`

	Organizations []Organization `gorm:"foreignKey:GovernorateID" json:"organizations,omitempty"`
}

func (Governorate) TableName() string {
	return "governorates"
}
