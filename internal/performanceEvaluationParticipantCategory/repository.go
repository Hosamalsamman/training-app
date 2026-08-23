package performanceEvaluationParticipantCategory

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

func (r *Repository) GetAll() ([]models.PerformanceEvaluationParticipantCategory, error) {

	var pCats []models.PerformanceEvaluationParticipantCategory

	err := r.db.
		Preload("Client").
		Preload("ParticipantType").
		Preload("PerformanceEvaluationCategory").
		Preload("Items").
		Find(&pCats).Error

	return pCats, err
}

func (r *Repository) GetByID(id int) (*models.PerformanceEvaluationParticipantCategory, error) {

	var pItem models.PerformanceEvaluationParticipantCategory

	err := r.db.
		Preload("Client").
		Preload("ParticipantType").
		Preload("PerformanceEvaluationCategory").
		Preload("Items").
		First(&pItem, id).Error

	if err != nil {
		return nil, err
	}

	return &pItem, nil
}