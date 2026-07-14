package organizations

import "training-app/db"

type Repository struct{}

func (r Repository) GetAll() ([]Organization, error) {

	var orgs []Organization

	err := db.DB.
		Preload("OrganizationType").
		Find(&orgs).Error

	return orgs, err
}

func (r Repository) GetByID(id string) (*Organization, error) {

	var org Organization

	err := db.DB.
		Preload("OrganizationType").
		First(&org, id).Error

	if err != nil {
		return nil, err
	}

	return &org, nil
}
