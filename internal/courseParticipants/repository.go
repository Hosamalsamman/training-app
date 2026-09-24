package courseParticipants

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

func (r *Repository) GetAll() ([]models.CourseParticipant, error) {

	var courseParticipants []models.CourseParticipant

	err := r.db.
	Preload("Course").
    Preload("Person").
    Preload("Client").
	Find(&courseParticipants).Error

	return courseParticipants, err
}

func (r *Repository) GetByID(id int) (*models.CourseParticipant, error) {

	var courseParticipant models.CourseParticipant

	err := r.db.
	Preload("Course").
    Preload("Person").
    Preload("Client").
	First(&courseParticipant, id).Error

	if err != nil {
		return nil, err
	}

	return &courseParticipant, nil
}

// Create performs the INSERT only. It does not commit: it is
// always called on a repository scoped to a transaction DB
// (WithDB) so the caller decides the commit.
func (r *Repository) Create(participant *models.CourseParticipant) error {
	return r.db.Create(participant).Error
}

// Update performs a full UPDATE of the mutable participant
// fields. It is always called on a client-scoped repository,
// so the WHERE carries the client_id and the update cannot
// cross tenants even if the id of another client's participant
// was sent.
func (r *Repository) Update(id int, participant *models.CourseParticipant) (int64, error) {

	res := r.db.Model(&models.CourseParticipant{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"course_id": participant.CourseID,
			"person_id": participant.PersonID,
		})

	return res.RowsAffected, res.Error
}

// Delete removes the participant by id. Like Update it is
// always called on a client-scoped repository, so participants
// of other clients are invisible to the WHERE clause.
func (r *Repository) Delete(id int) (int64, error) {

	res := r.db.Delete(&models.CourseParticipant{}, id)

	return res.RowsAffected, res.Error
}