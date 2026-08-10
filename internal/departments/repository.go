package departments

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

func (r *Repository) GetAll() ([]models.Department, error) {

	var departments []models.Department

	err := r.db.
		Preload("Client").
		Preload("Organization").
		Find(&departments).Error

	return departments, err
}

func (r *Repository) GetByID(id int) (*models.Department, error) {

	var department models.Department

	err := r.db.
		Preload("Client").
		Preload("Organization").
		First(&department, id).Error

	if err != nil {
		return nil, err
	}

	return &department, nil
}