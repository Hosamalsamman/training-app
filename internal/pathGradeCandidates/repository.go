package pathGradeCandidate

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

func (r *Repository) GetAll() ([]models.PathGradeCandidate, error) {

	var candidates []models.PathGradeCandidate

	err := r.db.
		Preload("Client").
		Preload("Person").
		Preload("PathGrade").
		Preload("PathGrade.Grade").
		Preload("PathGrade.LearningPath").
		Preload("Term").
		Preload("FinalEvaluation").
		Find(&candidates).Error

	return candidates, err
}

func (r *Repository) GetByID(id int) (*models.PathGradeCandidate, error) {

	var candidate models.PathGradeCandidate

	err := r.db.
		Preload("Client").
		Preload("Person").
		Preload("PathGrade").
		Preload("PathGrade.Grade").
		Preload("PathGrade.LearningPath").
		Preload("Term").
		Preload("FinalEvaluation").
		First(&candidate, id).Error

	if err != nil {
		return nil, err
	}

	return &candidate, nil
}

// Create performs the INSERT only. It does not commit: it is
// always called on a repository scoped to a transaction DB
// (WithDB) so the caller decides the commit.
func (r *Repository) Create(candidate *models.PathGradeCandidate) error {
	return r.db.Create(candidate).Error
}

// Update performs a full UPDATE of the mutable candidate
// fields. It is always called on a client-scoped repository,
// so the WHERE carries the client_id and the update cannot
// cross tenants even if the id of another client's candidate
// was sent. The final evaluation foreign key column is named
// final_evaluation in the database, without the _id suffix.
func (r *Repository) Update(id int, candidate *models.PathGradeCandidate) (int64, error) {

	res := r.db.Model(&models.PathGradeCandidate{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"person_id":        candidate.PersonID,
			"path_grade_id":    candidate.PathGradeID,
			"term_id":          candidate.TermID,
			"final_evaluation": candidate.FinalEvaluationID,
		})

	return res.RowsAffected, res.Error
}

// Delete removes the candidate by id. Like Update it is always
// called on a client-scoped repository, so candidates of other
// clients are invisible to the WHERE clause.
func (r *Repository) Delete(id int) (int64, error) {

	res := r.db.Delete(&models.PathGradeCandidate{}, id)

	return res.RowsAffected, res.Error
}
