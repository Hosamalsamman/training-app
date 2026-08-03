package learningSubjects

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

func (r *Repository) GetAll() ([]models.LearningSubject, error) {

	var learningSubjects []models.LearningSubject

	err := r.db.
		Preload("Client").
		Find(&learningSubjects).Error

	return learningSubjects, err
}

func (r *Repository) GetByID(id int) (*models.LearningSubject, error) {

	var learningSubject models.LearningSubject

	err := r.db.
		Preload("Client").
		First(&learningSubject, id).Error

	if err != nil {
		return nil, err
	}

	return &learningSubject, nil
}