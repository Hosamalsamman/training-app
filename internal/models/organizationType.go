package models

type OrganizationType struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	ClientID int `gorm:"column:client_id;not null" json:"client_id"`

	Organizations []Organization `gorm:"foreignKey:OrganizationTypeID" json:"organizations,omitempty"`
}

func (OrganizationType) TableName() string {
	return "organization_types"
}