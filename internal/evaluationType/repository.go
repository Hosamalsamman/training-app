package evaluationType

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

func (r *Repository) DB() *gorm.DB {
	return r.db
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

func (r *Repository) GetAll() ([]models.EvaluationType, error) {
	var eTypes []models.EvaluationType

	err := r.db.
		Preload("Client").
		Find(&eTypes).Error

	return eTypes, err
}

func (r *Repository) GetByID(id int) (*models.EvaluationType, error) {
	var eType models.EvaluationType

	err := r.db.
		Preload("Client").
		First(&eType, id).Error

	if err != nil {
		return nil, err
	}

	return &eType, nil
}
