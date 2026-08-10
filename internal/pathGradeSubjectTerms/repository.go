package pathGradeSubjectTerm

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

func (r *Repository) ForClient(clientID int) *Repository {
	return &Repository{
		db: r.db.Where("client_id = ?", clientID),
	}
}

func (r *Repository) GetAll() ([]models.PathGradeSubjectTerm, error) {

	var terms []models.PathGradeSubjectTerm

	err := r.db.
		Preload("Client").
		Preload("PathGradeSubject").
		Preload("Term").
		Preload("Documentations").
		Preload("PathGradeSubject.LearningPath").
		Find(&terms).Error

	return terms, err
}

func (r *Repository) GetByID(id int) (*models.PathGradeSubjectTerm, error) {

	var term models.PathGradeSubjectTerm

	err := r.db.
		Preload("Client").
		Preload("PathGradeSubject").
		Preload("Term").
		Preload("Documentations").
		First(&term, id).Error

	if err != nil {
		return nil, err
	}

	return &term, nil
}
