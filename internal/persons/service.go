package persons

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

func (s *Service) GetAll(clientID int) ([]models.Person, error) {

	repo := s.repo.ForClient(clientID)

	return repo.GetAll()
}

func (s *Service) GetByID(clientID int, id int) (*models.Person, error) {

	repo := s.repo.ForClient(clientID)

	return repo.GetByID(id)
}

func (s *Service) Create(
	clientID int,
	req *models.CreatePersonRequest,
) (*models.Person, error) {

	var person models.Person

	err := dbutil.WithTransaction(
		s.repo.DB(),
		func(tx *gorm.DB) error {

			// Use the transaction DB and restrict this repository
			// to the authenticated client's data.
			repo := s.repo.
				WithDB(tx).
				ForClient(clientID)

			// Build the database model from the request.
			person = models.Person{
				Code:                            req.Code,
				Name:                            req.Name,
				GovernorateID:                   req.GovernorateID,
				Address:                         req.Address,
				TelephoneWhatsapp:               req.TelephoneWhatsapp,
				QualificationID:                 req.QualificationID,
				JobID:                           req.JobID,
				OrganizationID:                  req.OrganizationID,
				DepartmentID:                    req.DepartmentID,
				WorkSiteID:                      req.WorkSiteID,
				JobTypeGroupID:                  req.JobTypeGroupID,
				CurrentLearningPathID:           req.CurrentLearningPathID,
				CurrentGradeID:                  req.CurrentGradeID,
				ContractDate:                    req.ContractDate,
				IsActive:                        true,
				IsTrainer:                       req.IsTrainer,
				TrainerCertifyingOrganizationID: req.TrainerCertifyingOrganizationID,

				// IMPORTANT:
				// client_id comes from JWT, not the request.
				ClientID: clientID,
			}

			// Repository only performs INSERT.
			// It does not commit.
			return repo.Create(&person)
		},
	)

	if err != nil {
		return nil, err
	}

	return &person, nil
}

// TODO: register person by giving him username and password and other required fields
