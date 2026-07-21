package organizations

import (
	"training-app/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) ForClient(clientID int) *Repository {

	return &Repository{
		db: r.db.Where("client_id = ?", clientID),
	}

}

func (r *Repository) GetAll() ([]models.Organization, error) {

	var orgs []models.Organization

	err := r.db.
    Preload("Client").
    Preload("OrganizationType").
    Preload("Governorate").
    Find(&orgs).Error

	return orgs, err
}

func (r *Repository) GetByID(id int) (*models.Organization, error) {

	var org models.Organization

	err := r.db.
	Preload("Client").
    Preload("OrganizationType").
    Preload("Governorate").
	First(&org, id).Error

	if err != nil {
		return nil, err
	}

	return &org, nil
}
