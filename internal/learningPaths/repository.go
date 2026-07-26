package learningPaths

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

func (r *Repository) GetAll() ([]models.LearningPath, error) {

	var learningPaths []models.LearningPath

	err := r.db.
		Preload("Client").
		Find(&learningPaths).Error

	return learningPaths, err
}

func (r *Repository) GetByID(id int) (*models.LearningPath, error) {

	var learningPath models.LearningPath

	err := r.db.
		Preload("Client").
		First(&learningPath, id).Error

	if err != nil {
		return nil, err
	}

	return &learningPath, nil
}