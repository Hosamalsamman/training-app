package jobs

import (
	"training-app/internal/dbutil"
	"training-app/internal/models"

	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(clientID int) ([]models.Job, error) {

	repo := s.repo.ForClient(clientID)

	return repo.GetAll()
}

func (s *Service) GetByID(clientID int, id int) (*models.Job, error) {

	repo := s.repo.ForClient(clientID)

	return repo.GetByID(id)
}

func (s *Service) Create(clientID int, req *models.Job) (*models.Job, error) {

	var job models.Job

	err := dbutil.WithTransaction(s.repo.DB(), func(tx *gorm.DB) error {

		repo := s.repo.
			WithDB(tx).
			ForClient(clientID)

		job = models.Job{
			Code:     req.Code,
			Name:     req.Name,
			ClientID: clientID,
		}

		return repo.Create(&job)
	})

	if err != nil {
		return nil, err
	}

	return &job, nil
}