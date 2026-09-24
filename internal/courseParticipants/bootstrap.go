package courseParticipants

import (
	"training-app/internal/courses"
	"training-app/internal/pathGradeSubjects"
	"training-app/internal/persons"

	"gorm.io/gorm"
)

func New(db *gorm.DB) *Handler {

	repo := NewRepository(db)

	// The courses, persons and pathGradeSubjects modules own
	// their tables: the course and person a participant
	// references, and the path grade of the course's subject,
	// are loaded through their repositories, including the
	// path-grade candidate rule of participants, instead of
	// this module re-implementing those lookups.
	courseRepo := courses.NewRepository(db)
	personRepo := persons.NewRepository(db)
	pathGradeSubjectRepo := pathGradeSubjects.NewRepository(db)

	service := NewService(repo, courseRepo, personRepo, pathGradeSubjectRepo)

	handler := NewHandler(service)

	return handler
}