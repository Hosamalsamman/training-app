package models

import "time"

type Person struct {
	ID   int     `gorm:"primaryKey;column:id" json:"id"`
	Code *string `gorm:"column:code;size:10" json:"code"`

	Name string `gorm:"column:name;size:100;not null" json:"name"`

	GovernorateID int         `gorm:"column:governorate_id;not null" json:"governorate_id"`
	Governorate   Governorate `gorm:"foreignKey:GovernorateID" json:"governorate"`

	Address           *string `gorm:"column:address;size:1500" json:"address"`
	TelephoneWhatsapp *string `gorm:"column:telephone_whatsapp;size:20" json:"telephone_whatsapp"`

	QualificationID int           `gorm:"column:qualification_id;not null" json:"qualification_id"`
	Qualification   Qualification `gorm:"foreignKey:QualificationID" json:"qualification"`

	JobID int `gorm:"column:job_id;not null" json:"job_id"`
	Job   Job `gorm:"foreignKey:JobID" json:"job"`

	OrganizationID int          `gorm:"column:organization_id;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"organization"`

	DepartmentID int        `gorm:"column:department_id;not null" json:"department_id"`
	Department   Department `gorm:"foreignKey:DepartmentID" json:"department"`

	WorkSiteID *int      `gorm:"column:work_site_id" json:"work_site_id"`
	WorkSite   *WorkSite `gorm:"foreignKey:WorkSiteID" json:"work_site"`

	JobTypeGroupID *int          `gorm:"column:job_type_group" json:"job_type_group"`
	JobTypeGroup   *JobTypeGroup `gorm:"foreignKey:JobTypeGroupID" json:"job_type_group_data"`

	CurrentLearningPathID int          `gorm:"column:current_learning_path_id;not null" json:"current_learning_path_id"`
	CurrentLearningPath   LearningPath `gorm:"foreignKey:CurrentLearningPathID" json:"current_learning_path"`

	CurrentGradeID int   `gorm:"column:current_grade;not null" json:"current_grade"`
	CurrentGrade   Grade `gorm:"foreignKey:CurrentGradeID" json:"current_grade_data"`

	DateOfCurrentGrade *time.Time `gorm:"column:date_ofcurrentgarde" json:"date_of_current_grade"`
	ContractDate       *time.Time `gorm:"column:contract_date" json:"contract_date"`

	IsActive  bool `gorm:"column:is_active;not null" json:"is_active"`
	IsTrainer bool `gorm:"column:is_trainer;not null" json:"is_trainer"`

	TrainerCertifyingOrganizationID *int          `gorm:"column:trainer_certifying_organization_id" json:"trainer_certifying_organization_id"`
	TrainerCertifyingOrganization   *Organization `gorm:"foreignKey:TrainerCertifyingOrganizationID" json:"trainer_certifying_organization"`

	UserName     *string `gorm:"column:user_name;size:500" json:"user_name"`
	UserPassword *string `gorm:"column:user_password;size:500" json:"user_password"`

	ClientID int    `gorm:"column:client_id;not null" json:"client_id"`
	Client   Client `gorm:"foreignKey:ClientID" json:"client"`

	TrainerSubjects []TrainerSubject `gorm:"foreignKey:PersonID" json:"trainer_subjects,omitempty"`
}

func (Person) TableName() string {
	return "persons"
}
