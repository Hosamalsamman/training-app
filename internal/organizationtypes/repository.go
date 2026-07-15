package organizationtypes

import "training-app/db"

type Repository struct{}

func (r Repository) GetAll() ([]OrganizationType, error) {
	var types []OrganizationType

	err := db.DB.Find(&types).Error

	return types, err
}

func (r Repository) GetByID(id int) (*OrganizationType, error) {
	var t OrganizationType

	err := db.DB.First(&t, id).Error

	if err != nil {
		return nil, err
	}

	return &t, nil
}
