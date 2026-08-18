package pathGradeSubjects

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

func (r *Repository) GetAll() ([]models.PathGradeSubject, error) {

	var pathGradeSubjects []models.PathGradeSubject

	err := r.db.
		Preload("Client").
		Preload("PathGrade").
		Preload("LearningSubject").
		Preload("PathGradeSubjectTerms.LearningTerm").
		Preload("Documentations").
		Find(&pathGradeSubjects).Error

	return pathGradeSubjects, err
}

func (r *Repository) GetByID(id int) (*models.PathGradeSubject, error) {

	var pathGradeSubject models.PathGradeSubject

	err := r.db.
		Preload("Client").
		Preload("PathGrade").
		Preload("LearningSubject").
		Preload("PathGradeSubjectTerms.LearningTerm").
		Preload("Documentations").
		First(&pathGradeSubject, id).Error

	if err != nil {
		return nil, err
	}

	return &pathGradeSubject, nil
}