package courseSessionParticipants

import (
	// "training-app/internal/dbutil"
	"training-app/internal/models"

	// "gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(clientID int) ([]models.CourseSessionParticipant, error) {

	repo := s.repo.ForClient(clientID)

	return repo.GetAll()
}

func (s *Service) GetByID(clientID int, id int) (*models.CourseSessionParticipant, error) {

	repo := s.repo.ForClient(clientID)

	return repo.GetByID(id)
}