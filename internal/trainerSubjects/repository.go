package trainerSubjects

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

func (r *Repository) GetAll() ([]models.TrainerSubject, error) {

	var trainerSubjects []models.TrainerSubject

	err := r.db.
		Preload("Client").
		Preload("Person").
		Preload("LearningSubject").
		Find(&trainerSubjects).Error

	return trainerSubjects, err
}

func (r *Repository) GetByID(id int) (*models.TrainerSubject, error) {

	var trainerSubject models.TrainerSubject

	err := r.db.
		Preload("Client").
		Preload("Person").
		Preload("LearningSubject").
		First(&trainerSubject, id).Error

	if err != nil {
		return nil, err
	}

	return &trainerSubject, nil
}