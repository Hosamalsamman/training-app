package courseSessionParticipants

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

func (r *Repository) GetAll() ([]models.CourseSessionParticipant, error) {

	var courseSessionParticipants []models.CourseSessionParticipant

	err := r.db.
	Preload("Client").
	Preload("Person").
	Preload("CourseSession").
	Preload("CourseSession.Course").
	Find(&courseSessionParticipants).Error

	return courseSessionParticipants, err
}

func (r *Repository) GetByID(id int) (*models.CourseSessionParticipant, error) {

	var courseSessionParticipant models.CourseSessionParticipant

	err := r.db.
	Preload("Client").
	Preload("Person").
	Preload("CourseSession").
	Preload("CourseSession.Course").
	First(&courseSessionParticipant, id).Error

	if err != nil {
		return nil, err
	}

	return &courseSessionParticipant, nil
}