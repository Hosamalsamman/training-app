package organizationtypes

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

func (r *Repository) GetAll() ([]models.OrganizationType, error) {
	var types []models.OrganizationType

	err := r.db.Find(&types).Error

	return types, err
}

func (r *Repository) GetByID(id int) (*models.OrganizationType, error) {
	var t models.OrganizationType

	err := r.db.First(&t, id).Error

	if err != nil {
		return nil, err
	}

	return &t, nil
}
