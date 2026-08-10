package courseParticipants

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

func (r *Repository) GetAll() ([]models.CourseParticipant, error) {

	var courseParticipants []models.CourseParticipant

	err := r.db.
	Preload("Course").
    Preload("Person").
    Preload("FinalEvaluation").
    Preload("Client").
	Find(&courseParticipants).Error

	return courseParticipants, err
}

func (r *Repository) GetByID(id int) (*models.CourseParticipant, error) {

	var courseParticipant models.CourseParticipant

	err := r.db.
	Preload("Course").
    Preload("Person").
    Preload("FinalEvaluation").
    Preload("Client").
	First(&courseParticipant, id).Error

	if err != nil {
		return nil, err
	}

	return &courseParticipant, nil
}