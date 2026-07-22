package jobTypeGroups

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

func (r *Repository) GetAll() ([]models.JobTypeGroup, error) {

	var groups []models.JobTypeGroup

	err := r.db.
		Preload("WorkGroup").
		Find(&groups).Error

	return groups, err
}

func (r *Repository) GetByID(id int) (*models.JobTypeGroup, error) {

	var group models.JobTypeGroup

	err := r.db.
		Preload("WorkGroup").
		First(&group, id).Error

	if err != nil {
		return nil, err
	}

	return &group, nil
}