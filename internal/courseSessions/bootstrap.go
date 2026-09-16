package courseSessions

import (
	"training-app/internal/courses"

	"gorm.io/gorm"
)

func New(db *gorm.DB) *Handler {

	repo := NewRepository(db)

	// The courses module owns its table: the course a
	// session belongs to is loaded through its repository,
	// including the executed-only rule of sessions, instead
	// of this module re-implementing a course lookup.
	courseRepo := courses.NewRepository(db)

	service := NewService(repo, courseRepo)

	handler := NewHandler(service)

	return handler
}