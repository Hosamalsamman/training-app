package models

type Department struct {
	ID   int     `gorm:"primaryKey;column:id" json:"id"`
	Code *string `gorm:"column:code;size:20" json:"code"`
	Name string  `gorm:"column:name;size:2000;not null" json:"name"`

	OrganizationID int          `gorm:"column:organization_id;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"organization"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`
}

func (Department) TableName() string {
	return "departments"
}