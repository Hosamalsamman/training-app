package models

type Group struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name;size:1000;not null" json:"name"`

	Persons []Person `gorm:"foreignKey:GroupID" json:"persons,omitempty"`
}

func (Group) TableName() string {
	return "groups"
}
