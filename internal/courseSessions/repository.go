package courseSessions

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

func (r *Repository) GetAll() ([]models.CourseSession, error) {
	var CourseSessions []models.CourseSession
	err := r.db.
	Preload("Client").
    Preload("Course").
	Preload("Documentations").
	Find(&CourseSessions).Error
	return CourseSessions, err
}

func (r *Repository) GetByID(id int) (*models.CourseSession, error) {
	var CourseSession models.CourseSession
	err := r.db.
	Preload("Client").
    Preload("Course").
	Preload("Documentations").
	First(&CourseSession, id).Error
	if err != nil {
		return nil, err
	}
	return &CourseSession, nil
}
