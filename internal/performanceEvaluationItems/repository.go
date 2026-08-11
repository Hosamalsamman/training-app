package performanceEvaluationItems

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

func (r *Repository) GetAll() ([]models.PerformanceEvaluationItem, error) {

	var pItems []models.PerformanceEvaluationItem

	err := r.db.
		Preload("Client").
		Find(&pItems).Error

	return pItems, err
}

func (r *Repository) GetByID(id int) (*models.PerformanceEvaluationItem, error) {

	var pItem models.PerformanceEvaluationItem

	err := r.db.
		Preload("Client").
		First(&pItem, id).Error

	if err != nil {
		return nil, err
	}

	return &pItem, nil
}