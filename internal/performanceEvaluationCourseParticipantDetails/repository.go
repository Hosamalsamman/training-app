package performanceEvaluationCourseParticipantDetails

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

func (r *Repository) GetAll() ([]models.PerformanceEvaluationCourseParticipantDetail, error) {

	var participantDetails []models.PerformanceEvaluationCourseParticipantDetail

	err := r.db.
		Preload("CourseParticipant.Course").
		Preload("CourseParticipant.Person").
		Preload("EnteredByPerson").
		Preload("Client").
		Preload("PerformanceEvaluationParticipantCategoryItem.PerformanceEvaluationParticipantCategory.ParticipantType").
		Preload("PerformanceEvaluationParticipantCategoryItem.PerformanceEvaluationParticipantCategory.PerformanceEvaluationCategory").
		Preload("PerformanceEvaluationParticipantCategoryItem.PerformanceEvaluationItem").
		Find(&participantDetails).Error

	return participantDetails, err
}

func (r *Repository) GetByID(id int) (*models.PerformanceEvaluationCourseParticipantDetail, error) {

	var participantDetail models.PerformanceEvaluationCourseParticipantDetail

	err := r.db.
		Preload("CourseParticipant.Course").
		Preload("CourseParticipant.Person").
		Preload("EnteredByPerson").
		Preload("Client").
		Preload("PerformanceEvaluationParticipantCategoryItem.PerformanceEvaluationParticipantCategory.ParticipantType").
		Preload("PerformanceEvaluationParticipantCategoryItem.PerformanceEvaluationParticipantCategory.PerformanceEvaluationCategory").
		Preload("PerformanceEvaluationParticipantCategoryItem.PerformanceEvaluationItem").
		First(&participantDetail, id).Error

	if err != nil {
		return nil, err
	}

	return &participantDetail, nil
}
