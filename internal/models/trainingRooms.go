package models


type TrainingRoom struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	WorkSiteID int      `gorm:"column:work_site_id;not null" json:"work_site_id"`
	WorkSite   WorkSite `gorm:"foreignKey:WorkSiteID" json:"work_site"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`

	Courses []Course `gorm:"foreignKey:RoomID" json:"courses,omitempty"`
}