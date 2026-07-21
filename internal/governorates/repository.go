package governorates

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

func (r *Repository) GetAll() ([]models.Governorate, error) {
	var govs []models.Governorate

	err := r.db.Find(&govs).Error

	return govs, err
}

func (r *Repository) GetByID(id int) (*models.Governorate, error) {
	var g models.Governorate

	err := r.db.First(&g, id).Error

	if err != nil {
		return nil, err
	}

	return &g, nil
}
