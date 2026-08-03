package learningTerms

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

func (r *Repository) GetAll() ([]models.LearningTerm, error) {

	var learningTerms []models.LearningTerm

	err := r.db.
		Preload("Client").
		Find(&learningTerms).Error

	return learningTerms, err
}

func (r *Repository) GetByID(id int) (*models.LearningTerm, error) {

	var learningTerm models.LearningTerm

	err := r.db.
		Preload("Client").
		First(&learningTerm, id).Error

	if err != nil {
		return nil, err
	}

	return &learningTerm, nil
}