package organizations

import "training-app/internal/organizationtypes"

// Organization — main table
type Organization struct {
	ID                 int                                `gorm:"primaryKey;column:id" json:"id"`
	Name               string                             `gorm:"column:name;size:2000;not null" json:"name"`
	OrganizationTypeID int                                `gorm:"column:organization_type_id;not null" json:"organization_type_id"`
	OrganizationType   organizationtypes.OrganizationType `gorm:"foreignKey:OrganizationTypeID" json:"organization_type"`
}
