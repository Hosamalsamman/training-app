package governorates

import "training-app/internal/models"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(clientID int) ([]models.Governorate, error) {
	repo := s.repo.ForClient(clientID)
	return repo.GetAll()
}

func (s *Service) GetByID(clientID int, id int) (*models.Governorate, error) {
	repo := s.repo.ForClient(clientID)
	return repo.GetByID(id)
}
