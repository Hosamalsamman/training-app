package courseSessionParticipants

import (
	"training-app/internal/courseParticipants"
	"training-app/internal/courseSessions"

	"gorm.io/gorm"
)

func New(db *gorm.DB) *Handler {

	repo := NewRepository(db)

	// The courseSessions and courseParticipants modules own
	// their tables: the session and the enrollment an
	// attendance row references are loaded through their
	// repositories, including the same-course rule, instead of
	// this module re-implementing those lookups.
	sessionRepo := courseSessions.NewRepository(db)
	participantRepo := courseParticipants.NewRepository(db)

	service := NewService(repo, sessionRepo, participantRepo)

	handler := NewHandler(service)

	return handler
}