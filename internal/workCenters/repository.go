package workCenters

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

func (r *Repository) forClient(clientID int) *Repository {

	return &Repository{
		db: r.db.Where("client_id = ?", clientID),
	}
}

func (r *Repository) GetAll() ([]models.WorkCenter, error) {

	var workCenters []models.WorkCenter

	err := r.db.
		Preload("Governorate").
		Preload("Organization").
		Find(&workCenters).Error

	return workCenters, err
}

func (r *Repository) GetByID(id int) (*models.WorkCenter, error) {

	var workCenter models.WorkCenter

	err := r.db.
		Preload("Governorate").
		Preload("Organization").
		First(&workCenter, id).Error

	if err != nil {
		return nil, err
	}

	return &workCenter, nil
}