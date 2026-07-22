package models

type Country struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	Name string `gorm:"column:name;size:100;not null" json:"name"`

	ClientID int `gorm:"column:client_id;not null" json:"client_id"`

	Governorates []Governorate `gorm:"foreignKey:CountryID" json:"governorates,omitempty"`
}

func (Country) TableName() string {
	return "countries"
}
