package countries

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

func (r *Repository) GetAll() ([]models.Country, error) {
	var countries []models.Country
	err := r.db.Preload("Client").Find(&countries).Error
	return countries, err
}

func (r *Repository) GetByID(id int) (*models.Country, error) {
	var country models.Country
	err := r.db.Preload("Client").First(&country, id).Error
	if err != nil {
		return nil, err
	}
	return &country, nil
}
