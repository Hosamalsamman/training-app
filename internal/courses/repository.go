package courses

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

func (r *Repository) GetAll() ([]models.Course, error) {
	var courses []models.Course
	err := r.db.
		Preload("Client").
		Preload("ExecutedCourses").
		Preload("PathGradeSubject").
		Preload("PathGradeSubjectTerm").
		Preload("LearningSubject").
		Preload("Room").
		Preload("FundingOrganization").
		Preload("Trainer").
		Preload("BackupTrainer").
		Preload("Coordinator").
		Preload("Sessions").
		Preload("Planned").
		Preload("Documentations").
		Find(&courses).Error
	return courses, err
}

func (r *Repository) GetByID(id int) (*models.Course, error) {
	var course models.Course
	err := r.db.
		Preload("Client").
		Preload("PathGradeSubject").
		Preload("PathGradeSubjectTerm").
		Preload("LearningSubject").
		Preload("Room").
		Preload("FundingOrganization").
		Preload("Trainer").
		Preload("BackupTrainer").
		Preload("Coordinator").
		Preload("Sessions").
		Preload("Planned").
		Preload("Documentations").
		First(&course, id).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

// notExecutedCondition is the SQL condition that expresses
// "this course has not been executed yet".
//
// Under the course lifecycle rules an execution is a NEW row
// (is_planned = true, is_executed = true) that points back at
// the planned course through planned_id, while the planned
// row itself keeps is_planned = true and is_executed = false
// forever. That means the flags on the planned row alone can
// never tell whether it was already executed: the only
// authoritative test is whether any row references it through
// planned_id. A planned course can be executed exactly once.
//
// The subquery needs no client_id filter: GetPlannableCourse
// only ever accepts a course of the caller's own client, so
// every planned_id value points at a same-client course, and
// the outer query is client-scoped anyway.
//
// NOT EXISTS is used instead of NOT IN because planned_id is
// nullable: a single NULL in a NOT IN subquery would make the
// whole NOT IN evaluate to no rows at all.
const notExecutedCondition = "NOT EXISTS (SELECT 1 FROM courses AS executions WHERE executions.planned_id = courses.id)"

// GetAllPlanned returns the client's pending planned courses
// (is_planned = true AND is_executed = false) that have not
// been executed yet, so the frontend only offers planned
// courses that can still be executed from. Every filter is
// optional: a nil filter means the query parameter was not
// sent, so no WHERE condition is added for it. This is
// GORM's chainable query, the equivalent of building an
// SQLAlchemy query with conditional .filter() calls.
func (r *Repository) GetAllPlanned(filters models.CourseListFilters) ([]models.Course, error) {

	var courses []models.Course

	// Base condition: only pending planned courses that were
	// not executed from yet.
	q := r.db.
		Where("is_planned = ? AND is_executed = ?", true, false).
		Where(notExecutedCondition)

	// Each filter is applied only when the frontend sent it.
	if filters.PathGradeSubjectID != nil {
		q = q.Where("path_grade_subject_id = ?", *filters.PathGradeSubjectID)
	}

	if filters.PathGradeSubjectTermID != nil {
		q = q.Where("path_grade_subject_term_id = ?", *filters.PathGradeSubjectTermID)
	}

	if filters.LearningSubjectID != nil {
		q = q.Where("learning_subject_id = ?", *filters.LearningSubjectID)
	}

	err := q.
		Preload("Client").
		Preload("PathGradeSubject").
		Preload("PathGradeSubjectTerm").
		Preload("LearningSubject").
		Preload("Room").
		Preload("FundingOrganization").
		Preload("Trainer").
		Preload("BackupTrainer").
		Preload("Coordinator").
		Preload("Sessions").
		Preload("Planned").
		Preload("Documentations").
		Find(&courses).Error

	return courses, err
}

// Create inserts a new course. It only performs the INSERT;
// the caller owns the transaction.
func (r *Repository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

// GetPlannableCourse fetches by id a course that can still
// be executed from. It is the backend safety net behind the
// frontend picker when creating an execution from a planned
// course.
//
// With the id in hand no other filter is needed beyond
// notExecutedCondition: the flags on the row itself cannot
// reveal whether it was already executed from (they never
// change), so the planned_id reference is the only
// authoritative test. The repository must be client-scoped,
// so courses of other clients are invisible here.
func (r *Repository) GetPlannableCourse(id int) (*models.Course, error) {

	var course models.Course

	err := r.db.
		Where(notExecutedCondition).
		First(&course, id).Error

	if err != nil {
		return nil, err
	}

	return &course, nil
}
