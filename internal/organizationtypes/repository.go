package organizationtypes

import "training-app/db"

type Repository struct{}

func (r Repository) GetAllOrganizationTypes() ([]OrganizationType, error) {
	var types []OrganizationType

	err := db.DB.Find(&types).Error

	return types, err
}

func (r Repository) GetOrganizationTypeByID(id string) (*OrganizationType, error) {
	var t OrganizationType

	err := db.DB.First(&t, id).Error

	if err != nil {
		return nil, err
	}

	return &t, nil
}
