package participantType

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

func (r *Repository) GetAll() ([]models.ParticipantType, error) {
	var types []models.ParticipantType

	err := r.db.
	Preload("Client").
	Find(&types).Error

	return types, err
}

func (r *Repository) GetByID(id int) (*models.ParticipantType, error) {
	var t models.ParticipantType

	err := r.db.
	Preload("Client").
	First(&t, id).Error

	if err != nil {
		return nil, err
	}

	return &t, nil
}