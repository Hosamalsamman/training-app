package clients

import (
	"training-app/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB // injected, not global
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ForClient(clientID int) *Repository {

	return &Repository{
		db: r.db.Where("client_id = ?", clientID),
	}

}

func (r *Repository) GetAll() ([]models.Client, error) {
	var clients []models.Client
	err := r.db.Find(&clients).Error
	return clients, err
}

func (r *Repository) GetByID(id int) (*models.Client, error) {
	var client models.Client
	err := r.db.First(&client, id).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}
