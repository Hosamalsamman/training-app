package performanceEvaluationCourseDetails

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

func (r *Repository) GetAll() ([]models.PerformanceEvaluationCourseDetail, error) {

	var courseDetails []models.PerformanceEvaluationCourseDetail

	err := r.db.
		Preload("Course").
		Preload("EnteredBy").
		Preload("Client").
		Preload("PerformanceEvaluationParticipantCategoryItem.PerformanceEvaluationParticipantCategory.ParticipantType").
		Preload("PerformanceEvaluationParticipantCategoryItem.PerformanceEvaluationParticipantCategory.PerformanceEvaluationCategory").
		Preload("PerformanceEvaluationParticipantCategoryItem.PerformanceEvaluationItem").
		Find(&courseDetails).Error

	return courseDetails, err
}

func (r *Repository) GetByID(id int) (*models.PerformanceEvaluationCourseDetail, error) {

	var courseDetail models.PerformanceEvaluationCourseDetail

	err := r.db.
		Preload("Course").
		Preload("EnteredBy").
		Preload("Client").
		Preload("PerformanceEvaluationParticipantCategoryItem.PerformanceEvaluationParticipantCategory.ParticipantType").
		Preload("PerformanceEvaluationParticipantCategoryItem.PerformanceEvaluationParticipantCategory.PerformanceEvaluationCategory").
		Preload("PerformanceEvaluationParticipantCategoryItem.PerformanceEvaluationItem").
		First(&courseDetail, id).Error

	if err != nil {
		return nil, err
	}

	return &courseDetail, nil
}
