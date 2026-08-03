package trainingRooms


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

func (r *Repository) GetAll() ([]models.TrainingRoom, error) {

	var TrainingRooms []models.TrainingRoom

	err := r.db.
	Preload("Client").
	Preload("WorkSite").
	Find(&TrainingRooms).Error

	return TrainingRooms, err
}

func (r *Repository) GetByID(id int) (*models.TrainingRoom, error) {

	var TrainingRoom models.TrainingRoom

	err := r.db.
	Preload("Client").
	Preload("WorkSite").
	Preload("Courses").
	First(&TrainingRoom, id).Error

	if err != nil {
		return nil, err
	}

	return &TrainingRoom, nil
}