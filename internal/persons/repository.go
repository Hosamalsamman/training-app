package persons

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

func (r *Repository) ForClient(clientID int) *Repository {
	return &Repository{
		db: r.db.Where("client_id = ?", clientID),
	}
}

func (r *Repository) GetAll() ([]models.Person, error) {

	var persons []models.Person

	err := r.db.
		Preload("Client").
		Preload("Governorate").
		Preload("Qualification").
		Preload("Job").
		Preload("Organization").
		Preload("Department").
		Preload("WorkSite").
		Preload("JobTypeGroup").
		Preload("CurrentLearningPath").
		Preload("CurrentGrade").
		Preload("TrainerCertifyingOrganization").
		Find(&persons).Error

	return persons, err
}

func (r *Repository) GetByID(id int) (*models.Person, error) {

	var person models.Person

	err := r.db.
		Preload("Client").
		Preload("Governorate").
		Preload("Qualification").
		Preload("Job").
		Preload("Organization").
		Preload("Department").
		Preload("WorkSite").
		Preload("JobTypeGroup").
		Preload("CurrentLearningPath").
		Preload("CurrentGrade").
		Preload("TrainerCertifyingOrganization").
		First(&person, id).Error

	if err != nil {
		return nil, err
	}

	return &person, nil
}