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
		Preload("PathGrade.Grade").
		Preload("LearningSubject").
		Preload("PathGradeSubjectTerms.Term").
		Preload("Documentations").
		Find(&pathGradeSubjects).Error

	return pathGradeSubjects, err
}

func (r *Repository) GetByID(id int) (*models.PathGradeSubject, error) {

	var pathGradeSubject models.PathGradeSubject

	err := r.db.
		Preload("Client").
		Preload("PathGrade").
		Preload("PathGrade.Grade").
		Preload("LearningSubject").
		Preload("PathGradeSubjectTerms.Term").
		Preload("Documentations").
		First(&pathGradeSubject, id).Error

	if err != nil {
		return nil, err
	}
	return &pathGradeSubject, nil
}

// GetByIDPlain fetches a path grade subject by id without any
// preloads. It is meant for validations that only need the
// foreign keys of the row (path_grade_id), inside a
// transaction. The repository must be client-scoped, so rows
// of other clients are invisible here.
func (r *Repository) GetByIDPlain(id int) (*models.PathGradeSubject, error) {

	var pathGradeSubject models.PathGradeSubject

	err := r.db.First(&pathGradeSubject, id).Error

	if err != nil {
		return nil, err
	}

	return &pathGradeSubject, nil
}