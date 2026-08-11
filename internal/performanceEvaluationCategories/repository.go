package performanceEvaluationCategories

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

func (r *Repository) ForClient(clientID int) *Repository {
	return &Repository{
		db: r.db.Where("client_id = ?", clientID),
	}
}

func (r *Repository) GetAll() ([]models.PerformanceEvaluationCategory, error) {

	var pCats []models.PerformanceEvaluationCategory

	err := r.db.
		Preload("Client").
		Find(&pCats).Error

	return pCats, err
}

func (r *Repository) GetByID(id int) (*models.PerformanceEvaluationCategory, error) {

	var pCat models.PerformanceEvaluationCategory

	err := r.db.
		Preload("Client").
		First(&pCat, id).Error

	if err != nil {
		return nil, err
	}

	return &pCat, nil
}