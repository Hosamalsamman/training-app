package pathGradeCandidate

import (
	"training-app/internal/pathGrades"
	"training-app/internal/persons"

	"gorm.io/gorm"
)

func New(db *gorm.DB) *Handler {

	repo := NewRepository(db)

	// The persons and pathGrades modules own their tables:
	// the candidate's person and target path grade are loaded
	// through their repositories, including the grade rules,
	// instead of this module re-implementing those lookups.
	personRepo := persons.NewRepository(db)
	pathGradeRepo := pathGrades.NewRepository(db)

	service := NewService(repo, personRepo, pathGradeRepo)

	handler := NewHandler(service)

	return handler
}
