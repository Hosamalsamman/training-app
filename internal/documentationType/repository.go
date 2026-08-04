package documentationType

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

func (r *Repository) WithDB(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) ForClient(clientID int) *Repository {
	return &Repository{
		db: r.db.Where("client_id = ?", clientID),
	}
}

func (r *Repository) GetAll() ([]models.DocumentationType, error) {
	var dTypes []models.DocumentationType

	err := r.db.
		Preload("Documentations").
		Find(&dTypes).Error

	return dTypes, err
}

func (r *Repository) GetByID(id int) (*models.DocumentationType, error) {
	var dType models.DocumentationType

	err := r.db.
		Preload("Documentations").
		First(&dType, id).Error

	if err != nil {
		return nil, err
	}

	return &dType, nil
}
