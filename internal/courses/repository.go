package courses

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

func (r *Repository) GetAll() ([]models.Course, error) {
	var courses []models.Course
	err := r.db.
		Preload("Client").
		Preload("ExecutedCourses").
		Preload("PathGradeSubjectTerm").
		Preload("LearningSubject").
		Preload("Room").
		Preload("FundingOrganization").
		Preload("Trainer").
		Preload("BackupTrainer").
		Preload("Coordinator").
		Preload("Sessions").
		Preload("Planned").Find(&courses).Error
	return courses, err
}

func (r *Repository) GetByID(id int) (*models.Course, error) {
	var course models.Course
	err := r.db.
		Preload("Client").
		Preload("PathGradeSubjectTerm").
		Preload("LearningSubject").
		Preload("Room").
		Preload("FundingOrganization").
		Preload("Trainer").
		Preload("BackupTrainer").
		Preload("Coordinator").
		Preload("Sessions").
		Preload("Planned").First(&course, id).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}
