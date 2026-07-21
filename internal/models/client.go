package models

type Client struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	Countries []Country `gorm:"foreignKey:ClientID" json:"countries,omitempty"`

	OrganizationTypes []OrganizationType `gorm:"foreignKey:ClientID" json:"organization_types,omitempty"`

	Organizations []Organization `gorm:"foreignKey:ClientID" json:"organizations,omitempty"`
}

func (Client) TableName() string {
	return "clients"
}