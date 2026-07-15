package countries

type Country struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name;size:100;not null" json:"name"`
}

func (Country) TableName() string {
	return "countries"
}
