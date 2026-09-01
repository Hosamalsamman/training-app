package persons

import "gorm.io/gorm"

func New(db *gorm.DB, secret string) *Handler {

	repo := NewRepository(db)

	service := NewService(repo)

	handler := NewHandler(service, secret)

	return handler
}
