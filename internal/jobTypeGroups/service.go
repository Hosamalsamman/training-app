package jobTypeGroups

import "training-app/internal/models"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(clientID int) ([]models.JobTypeGroup, error) {

	repo := s.repo.forClient(clientID)

	return repo.GetAll()
}

func (s *Service) GetByID(clientID, id int) (*models.JobTypeGroup, error) {

	repo := s.repo.forClient(clientID)

	return repo.GetByID(id)
}