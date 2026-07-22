package workGroups

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

func (r *Repository) GetAll() ([]models.WorkGroup, error) {

	var workGroups []models.WorkGroup

	err := r.db.
		Preload("JobTypeGroups").
		Find(&workGroups).Error

	return workGroups, err
}

func (r *Repository) GetByID(id int) (*models.WorkGroup, error) {

	var workGroup models.WorkGroup

	err := r.db.
		Preload("JobTypeGroups").
		First(&workGroup, id).Error

	if err != nil {
		return nil, err
	}

	return &workGroup, nil
}