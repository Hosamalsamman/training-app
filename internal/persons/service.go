package persons

import (
	"errors"
	"training-app/internal/dbutil"
	"training-app/internal/models"

	"golang.org/x/crypto/bcrypt"
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

// Sentinel errors for the register-user flow.
// The handler maps them to proper HTTP status codes.
var (
	ErrPersonNotFound    = errors.New("person not found")
	ErrAlreadyRegistered = errors.New("person already has credentials")
)

// RegisterUser turns an existing person (one that has no
// credentials yet) into a user by giving him a username,
// a password and a group.
func (s *Service) RegisterUser(
	clientID int,
	req *models.RegisterUserRequest,
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

			// The person must exist and belong to this client.
			existing, err := repo.GetByID(req.PersonID)

			if err != nil {

				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrPersonNotFound
				}

				return err
			}

			// A person can only become a user once.
			if existing.Username != nil && *existing.Username != "" {
				return ErrAlreadyRegistered
			}

			// Hash the password. The plain text is never stored.
			hashedPassword, err := bcrypt.GenerateFromPassword(
				[]byte(req.Password),
				bcrypt.DefaultCost,
			)

			if err != nil {
				return err
			}

			// Build the response model from the fetched person
			// and the new credentials.
			hashed := string(hashedPassword)

			person = *existing
			person.Username = &req.Username
			person.Password = &hashed
			person.GroupID = &req.GroupID

			// Repository only performs UPDATE.
			// It does not commit.
			return repo.UpdateCredentials(
				req.PersonID,
				req.Username,
				string(hashedPassword),
				req.GroupID,
			)
		},
	)

	if err != nil {
		return nil, err
	}

	return &person, nil
}
