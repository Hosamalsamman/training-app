package workSites

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

func (r *Repository) GetAll() ([]models.WorkSite, error) {

	var workSites []models.WorkSite

	err := r.db.
		Preload("WorkCenter").
		Find(&workSites).Error

	return workSites, err
}

func (r *Repository) GetByID(id int) (*models.WorkSite, error) {

	var workSite models.WorkSite

	err := r.db.
		Preload("WorkCenter").
		First(&workSite, id).Error

	if err != nil {
		return nil, err
	}

	return &workSite, nil
}