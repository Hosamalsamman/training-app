package groups

import (
	"errors"
	"training-app/internal/models"
	"training-app/internal/persons"

	"gorm.io/gorm"
)

type Service struct {
	repo       *Repository
	personRepo *persons.Repository
}

func NewService(repo *Repository, personRepo *persons.Repository) *Service {
	return &Service{
		repo:       repo,
		personRepo: personRepo,
	}
}

func (s *Service) GetAll() ([]models.Group, error) {

	return s.repo.GetAll()
}

func (s *Service) GetByID(id int) (*models.Group, error) {

	return s.repo.GetByID(id)
}

// Sentinel errors for the allowed-groups flow.
// The handler maps them to proper HTTP status codes.
var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserHasNoGroup = errors.New("user has no group assigned")
)

// GetAllowed returns the groups the given user is allowed to
// assign other users to: every group with an id greater than
// or equal to the user's own group id.
func (s *Service) GetAllowed(clientID int, userID int) ([]models.Group, error) {

	// The user must exist and belong to this client.
	groupID, err := s.personRepo.
		ForClient(clientID).
		GetGroupID(userID)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	// A user without a group cannot assign anyone.
	if groupID == nil {
		return nil, ErrUserHasNoGroup
	}

	return s.repo.GetAllowed(*groupID)
}
