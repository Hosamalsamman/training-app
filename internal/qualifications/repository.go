package qualifications

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

func (r *Repository) GetAll() ([]models.Qualification, error) {

	var qualifications []models.Qualification

	err := r.db.
		Preload("QualificationType").
		Find(&qualifications).Error

	return qualifications, err
}

func (r *Repository) GetByID(id int) (*models.Qualification, error) {

	var qualification models.Qualification

	err := r.db.
		Preload("QualificationType").
		First(&qualification, id).Error

	if err != nil {
		return nil, err
	}

	return &qualification, nil
}
