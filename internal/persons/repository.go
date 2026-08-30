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
		Preload("TrainerSubjects").
		Preload("TrainerCourses").
		Preload("BackupTrainerCourses").
		Preload("CoordinatedCourses").
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
		Preload("TrainerSubjects.LearningSubject").
		Preload("TrainerCourses.Room").
		Preload("TrainerCourses.LearningSubject").
		Preload("TrainerCourses.PathGradeSubjectTerm").
		Preload("BackupTrainerCourses.Room").
		Preload("BackupTrainerCourses.LearningSubject").
		Preload("CoordinatedCourses.Room").
		Preload("CoordinatedCourses.LearningSubject").
		First(&person, id).Error

	if err != nil {
		return nil, err
	}

	return &person, nil
}

func (r *Repository) Create(person *models.Person) error {
	return r.db.Create(person).Error
}

// UpdateCredentials sets the username, hashed password and
// group of an existing person. It is called on a client-scoped
// repository so the update cannot cross tenants.
func (r *Repository) UpdateCredentials(
	personID int,
	username string,
	hashedPassword string,
	groupID int,
) error {

	return r.db.Model(&models.Person{}).
		Where("id = ?", personID).
		Updates(map[string]any{
			"username": username,
			"password": hashedPassword,
			"group_id": groupID,
		}).Error
}

// GetGroupID returns only the group of a person without
// loading any associations. It is used by the groups module
// to decide which groups a user may be assigned to. It is
// called on a client-scoped repository so the lookup cannot
// cross tenants.
func (r *Repository) GetGroupID(personID int) (*int, error) {

	var person models.Person

	err := r.db.
		Select("id", "group_id").
		First(&person, personID).Error

	if err != nil {
		return nil, err
	}

	return person.GroupID, nil
}
