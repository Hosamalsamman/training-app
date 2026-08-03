package jobs

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

func (r *Repository) GetAll() ([]models.Job, error) {

	var jobs []models.Job

	err := r.db.Preload("Client").Find(&jobs).Error

	return jobs, err
}

func (r *Repository) GetByID(id int) (*models.Job, error) {

	var job models.Job

	err := r.db.Preload("Client").First(&job, id).Error

	if err != nil {
		return nil, err
	}

	return &job, nil
}