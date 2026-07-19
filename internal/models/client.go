package models

type Client struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name;size:2000;not null" json:"name"`
}

func (Client) TableName() string {
	return "clients"
}