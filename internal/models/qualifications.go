package models

type Qualification struct {
	ID int `gorm:"primaryKey;column:id" json:"id"`

	Code string `gorm:"column:code;size:10" json:"code,omitempty"`

	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	QualificationTypeID int `gorm:"column:qualification_type_id;not null" json:"qualification_type_id"`
	QualificationType QualificationType `gorm:"foreignKey:QualificationTypeID" json:"qualification_type"`

	ClientID int `gorm:"column:client_id;not null" json:"client_id"`
}

func (Qualification) TableName() string {
	return "qualifications"
}