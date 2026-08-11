package courseParticipantFinalExam

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

func (r *Repository) GetAll() ([]models.CourseParticipantFinalExam, error) {

	var participantsExams []models.CourseParticipantFinalExam

	err := r.db.
	Preload("Client").
	Preload("Organizer").
	Preload("CourseParticipant").
	Preload("Documentation").
	Preload("CourseParticipant.Person").
	Preload("CourseParticipant.Course").
	Find(&participantsExams).Error

	return participantsExams, err
}

func (r *Repository) GetByID(id int) (*models.CourseParticipantFinalExam, error) {

	var courseParticipantExam models.CourseParticipantFinalExam

	err := r.db.
	Preload("Client").
	Preload("Organizer").
	Preload("CourseParticipant").
	Preload("Documentation").
	Preload("CourseParticipant.Person").
	Preload("CourseParticipant.Course").
	First(&courseParticipantExam, id).Error

	if err != nil {
		return nil, err
	}

	return &courseParticipantExam, nil
}