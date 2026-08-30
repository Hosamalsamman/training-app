package groups

import (
	"training-app/internal/persons"

	"gorm.io/gorm"
)

func New(db *gorm.DB) *Handler {

	repo := NewRepository(db)

	// Needed to look up the group of the authenticated
	// user when listing the allowed groups.
	personRepo := persons.NewRepository(db)

	service := NewService(repo, personRepo)

	handler := NewHandler(service)

	return handler
}
