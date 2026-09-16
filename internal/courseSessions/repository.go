package courseSessions

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

func (r *Repository) GetAll() ([]models.CourseSession, error) {
	var CourseSessions []models.CourseSession
	err := r.db.
		Preload("Client").
		Preload("Course").
		Preload("Documentations").
		Find(&CourseSessions).Error
	return CourseSessions, err
}

func (r *Repository) GetByID(id int) (*models.CourseSession, error) {
	var CourseSession models.CourseSession
	err := r.db.
		Preload("Client").
		Preload("Course").
		Preload("Documentations").
		First(&CourseSession, id).Error
	if err != nil {
		return nil, err
	}
	return &CourseSession, nil
}

// Create performs the INSERT only. It does not commit: it is
// always called on a repository scoped to a transaction DB
// (WithDB) so the caller decides the commit.
func (r *Repository) Create(session *models.CourseSession) error {
	return r.db.Create(session).Error
}

// Update performs a full UPDATE of the mutable session fields.
// It is always called on a client-scoped repository, so the
// WHERE carries the client_id and the update cannot cross
// tenants even if the id of another client's session was sent.
func (r *Repository) Update(id int, session *models.CourseSession) (int64, error) {

	res := r.db.Model(&models.CourseSession{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"course_id":    session.CourseID,
			"session_date": session.SessionDate,
			"start_time":   session.StartTime,
			"end_time":     session.EndTime,
		})

	return res.RowsAffected, res.Error
}

// Delete removes the session by id. Like Update it is always
// called on a client-scoped repository, so sessions of other
// clients are invisible to the WHERE clause.
func (r *Repository) Delete(id int) (int64, error) {

	res := r.db.Delete(&models.CourseSession{}, id)

	return res.RowsAffected, res.Error
}
