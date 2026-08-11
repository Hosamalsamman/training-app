package documentations

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

func (r *Repository) GetAll() ([]models.Documentation, error) {

	var docs []models.Documentation

	err := r.db.
		Preload("Client").
		Preload("DocumentationType").
		Preload("LearningSubject").
		Preload("PathGradeSubject").
		Preload("PathGradeSubjectTerm").
		Preload("Course").
		Preload("CourseSession").
		Preload("CourseParticipantFinalExams").
		Find(&docs).Error

	return docs, err
}

func (r *Repository) GetByID(id int) (*models.Documentation, error) {

	var doc models.Documentation
	err := r.db.
		Preload("Client").
		Preload("DocumentationType").
		Preload("LearningSubject").
		Preload("PathGradeSubject").
		Preload("PathGradeSubjectTerm").
		Preload("Course").
		Preload("CourseSession").
		Preload("CourseParticipantFinalExams").
		First(&doc, id).Error

	if err != nil {
		return nil, err
	}

	return &doc, nil
}