package pathGrades

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

func (r *Repository) GetAll() ([]models.PathGrade, error) {
	var pGrades []models.PathGrade

	err := r.db.
	Preload("Client").
	Preload("LearningPath").
	Preload("Grade").
	Preload("PathGradeSubjects").
	Find(&pGrades).Error

	return pGrades, err
}

func (r *Repository) GetByID(id int) (*models.PathGrade, error) {
	var pGrade models.PathGrade

	err := r.db.
	Preload("Client").
	Preload("LearningPath").
	Preload("Grade").
	Preload("PathGradeSubjects").
	First(&pGrade, id).Error

	if err != nil {
		return nil, err
	}

	return &pGrade, nil
}