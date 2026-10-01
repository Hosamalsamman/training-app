package courseSessionParticipants

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

func (r *Repository) GetAll() ([]models.CourseSessionParticipant, error) {

	var courseSessionParticipants []models.CourseSessionParticipant

	err := r.db.
		Preload("Client").
		Preload("CourseParticipant").
		Preload("CourseParticipant.Person").
		Preload("CourseParticipant.Course").
		Preload("CourseSession").
		Preload("CourseSession.Course").
		Find(&courseSessionParticipants).Error

	return courseSessionParticipants, err
}

func (r *Repository) GetByID(id int) (*models.CourseSessionParticipant, error) {

	var courseSessionParticipant models.CourseSessionParticipant

	err := r.db.
		Preload("Client").
		Preload("CourseParticipant").
		Preload("CourseParticipant.Person").
		Preload("CourseParticipant.Course").
		Preload("CourseSession").
		Preload("CourseSession.Course").
		First(&courseSessionParticipant, id).Error

	if err != nil {
		return nil, err
	}

	return &courseSessionParticipant, nil
}

// Create performs the INSERT only. It does not commit: it is
// always called on a repository scoped to a transaction DB
// (WithDB) so the caller decides the commit.
func (r *Repository) Create(participant *models.CourseSessionParticipant) error {
	return r.db.Create(participant).Error
}

// Update performs a full UPDATE of the mutable fields. It is
// always called on a client-scoped repository, so the WHERE
// carries the client_id and the update cannot cross tenants
// even if the id of another client's row was sent.
func (r *Repository) Update(id int, participant *models.CourseSessionParticipant) (int64, error) {

	res := r.db.Model(&models.CourseSessionParticipant{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"course_session_id":     participant.CourseSessionID,
			"course_participant_id": participant.CourseParticipantID,
		})

	return res.RowsAffected, res.Error
}

// Delete removes the attendance row by id. Like Update it is
// always called on a client-scoped repository, so rows of
// other clients are invisible to the WHERE clause.
func (r *Repository) Delete(id int) (int64, error) {

	res := r.db.Delete(&models.CourseSessionParticipant{}, id)

	return res.RowsAffected, res.Error
}