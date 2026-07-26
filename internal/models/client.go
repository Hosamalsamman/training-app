package models

type Client struct {
	ID   int    `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"column:name;size:2000;not null" json:"name"`

	Countries         []Country          `gorm:"foreignKey:ClientID" json:"countries,omitempty"`
	Governorates      []Governorate      `gorm:"foreignKey:ClientID" json:"governorates,omitempty"` // if applicable
	OrganizationTypes []OrganizationType `gorm:"foreignKey:ClientID" json:"organization_types,omitempty"`
	Organizations     []Organization     `gorm:"foreignKey:ClientID" json:"organizations,omitempty"`
	Jobs              []Job              `gorm:"foreignKey:ClientID" json:"jobs,omitempty"`
	Grades            []Grade            `gorm:"foreignKey:ClientID" json:"grades,omitempty"`
	LearningPaths 	  []LearningPath `gorm:"foreignKey:ClientID" json:"learning_paths,omitempty"`
	Departments 	  []Department `gorm:"foreignKey:ClientID" json:"departments,omitempty"`
	TrainerSubjects   []TrainerSubject `gorm:"foreignKey:ClientID" json:"trainer_subjects,omitempty"`
}

func (Client) TableName() string {
	return "clients"
}