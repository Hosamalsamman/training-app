package workCenters

import "gorm.io/gorm"

func New(db *gorm.DB) *Handler {

	repo := NewRepository(db)

	service := NewService(repo)

	handler := NewHandler(service)

	return handler
}