package organizations

import "training-app/internal/models"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll() ([]models.Organization, error) {
	return s.repo.GetAll()
}

func (s *Service) GetByID(id int) (*models.Organization, error) {
	return s.repo.GetByID(id)
}
