package organizationtypes

type OrganizationType struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name;size:2000" json:"name"`
}
