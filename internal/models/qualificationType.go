package models

type QualificationType struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	Name string `gorm:"column:name;size:100;not null" json:"name"`

	ClientID int `gorm:"column:client_id;not null" json:"client_id"`

	Qualifications []Qualification `gorm:"foreignKey:QualificationTypeID" json:"qualifications,omitempty"`
}

func (QualificationType) TableName() string {
	return "qualification_types"
}

