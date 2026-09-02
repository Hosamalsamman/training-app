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

// GetByUsername returns the active person carrying the given
// username. Usernames are unique, so at most one account can
// match. It is NOT client-scoped: at login time the tenant is
// not known yet. Only the columns needed to verify credentials
// are selected. It is used by the login flow.
func (r *Repository) GetByUsername(username string) (*models.Person, error) {

	var person models.Person

	err := r.db.
		Select("id", "name", "client_id", "group_id", "password").
		Where("username = ? AND is_active = ?", username, true).
		First(&person).Error

	if err != nil {
		return nil, err
	}

	return &person, nil
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

// GetPasswordByID returns the password hash of a person
// without loading any associations. A nil password means the
// person was never registered as a user. It is called on a
// client-scoped repository so the lookup cannot cross
// tenants. It is used by the change-password and
// reset-password flows.
func (r *Repository) GetPasswordByID(personID int) (*string, error) {

	var person models.Person

	err := r.db.
		Select("id", "password").
		First(&person, personID).Error

	if err != nil {
		return nil, err
	}

	return person.Password, nil
}

// UpdatePassword sets only the password column of a person.
// It is called on a client-scoped repository so the update
// cannot cross tenants. It is used by the change-password and
// reset-password flows.
func (r *Repository) UpdatePassword(personID int, hashedPassword string) error {

	return r.db.Model(&models.Person{}).
		Where("id = ?", personID).
		Update("password", hashedPassword).Error
}

// GetAllUsers returns every person of the client who already
// has an account, i.e. every person the register-user flow
// would reject as already registered. The frontend renders
// this list so the admin can only pick a real user as a
// reset-password target. No associations are preloaded: the
// response is a small summary, not a full person. It is
// called on a client-scoped repository so the lookup cannot
// cross tenants.
func (r *Repository) GetAllUsers() ([]models.UserSummary, error) {

	var users []models.UserSummary

	// Model points at the persons table; the summary struct
	// is only the scan destination.
	err := r.db.
		Model(&models.Person{}).
		Select("id", "code", "name", "username", "group_id", "is_active").
		Where("username IS NOT NULL AND username <> ''").
		Find(&users).Error

	return users, err
}
