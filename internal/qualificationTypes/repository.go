package qualificationTypes

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

func (r *Repository) GetAll() ([]models.QualificationType, error) {

	var qualificationTypes []models.QualificationType

	err := r.db.Preload("Qualifications").Find(&qualificationTypes).Error

	return qualificationTypes, err
}

func (r *Repository) GetByID(id int) (*models.QualificationType, error) {

	var qualificationType models.QualificationType

	err := r.db.Preload("Qualifications").First(&qualificationType, id).Error
	if err != nil {
		return nil, err
	}

	return &qualificationType, nil
}