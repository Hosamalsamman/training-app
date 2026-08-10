package grades

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

func (r *Repository) GetAll() ([]models.Grade, error) {

	var grades []models.Grade

	err := r.db.Preload("Client").Find(&grades).Error

	return grades, err
}

func (r *Repository) GetByID(id int) (*models.Grade, error) {

	var grade models.Grade

	err := r.db.Preload("Client").First(&grade, id).Error

	if err != nil {
		return nil, err
	}

	return &grade, nil
}