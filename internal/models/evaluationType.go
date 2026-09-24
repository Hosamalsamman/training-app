package models

type EvaluationType struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	// Whether a candidate evaluated with this type must
	// repeat the course: should_repeat = 1 means the person
	// failed and may retake it, should_repeat = 0 means the
	// person passed (e.g. ناجح, ناجح بامتياز) and cannot
	// take the course again.
	ShouldRepeat *bool `gorm:"column:should_repeat" json:"should_repeat"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`

}