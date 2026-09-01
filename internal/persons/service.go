package persons

import (
	"errors"
	"strings"
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
	ErrPersonNotFound = errors.New("person not found")
	// The person already has a username and password.
	ErrAlreadyRegistered = errors.New("person already has credentials")
	// The authenticated user does not exist for this client.
	ErrUserNotFound = errors.New("authenticated user not found")
	// The authenticated user has no group, so he cannot assign anyone.
	ErrCallerHasNoGroup = errors.New("authenticated user has no group assigned")
	// The requested group is above the caller's own group.
	ErrGroupNotAllowed = errors.New("group is not allowed for this user")
	// The requested group does not exist.
	ErrGroupNotFound = errors.New("group does not exist")
	// The username became empty after normalization.
	ErrInvalidUsername = errors.New("username is empty")
	// Login failed: wrong username or wrong password.
	// Deliberately vague so the response cannot reveal
	// which part was wrong.
	ErrInvalidCredentials = errors.New("invalid username or password")
)

// RegisterUser turns an existing person (one that has no
// credentials yet) into a user by giving him a username,
// a password and a group.
//
// userID is the authenticated user performing the registration:
// he may only assign the new user to a group at or below his own
// group, the same rule behind the /allowed-groups route.
func (s *Service) RegisterUser(
	clientID int,
	userID int,
	req *models.RegisterUserRequest,
) (*models.Person, error) {

	// Normalize the username before it reaches the database:
	// trim surrounding whitespace and lowercase it.
	username := strings.ToLower(strings.TrimSpace(req.Username))

	if username == "" {
		return nil, ErrInvalidUsername
	}

	var person models.Person

	err := dbutil.WithTransaction(
		s.repo.DB(),
		func(tx *gorm.DB) error {

			// Use the transaction DB and restrict this repository
			// to the authenticated client's data.
			repo := s.repo.
				WithDB(tx).
				ForClient(clientID)

			// A user may only assign others to groups with an id
			// greater than or equal to his own group id, so the
			// group of the authenticated user is needed first.
			callerGroupID, err := repo.GetGroupID(userID)

			if err != nil {

				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrUserNotFound
				}

				return err
			}

			// A caller without a group cannot register anyone.
			if callerGroupID == nil {
				return ErrCallerHasNoGroup
			}

			// Reject any group above the caller's own group,
			// even if the frontend only offers allowed groups.
			if req.GroupID < *callerGroupID {
				return ErrGroupNotAllowed
			}

			// The requested group must exist. groups is a global
			// lookup table, so the check runs without client
			// scoping.
			exists, err := s.repo.WithDB(tx).GroupExists(req.GroupID)

			if err != nil {
				return err
			}

			if !exists {
				return ErrGroupNotFound
			}

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
			person.Username = &username
			person.Password = &hashed
			person.GroupID = &req.GroupID

			// Repository only performs UPDATE.
			// It does not commit.
			return repo.UpdateCredentials(
				req.PersonID,
				username,
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

// Login verifies the given credentials and returns the
// authenticated person on success.
//
// The username is normalized the same way the register-user
// flow normalizes it before storing it, so the user can type
// his username in any case.
func (s *Service) Login(
	username string,
	password string,
) (*models.Person, error) {

	// Same normalization as the register-user flow:
	// trim surrounding whitespace and lowercase it.
	username = strings.ToLower(strings.TrimSpace(username))

	if username == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	// Usernames are unique, so the first match is the
	// account to verify. An unknown username and a wrong
	// password both end up as invalid credentials, so the
	// response cannot reveal which part failed.
	person, err := s.repo.GetByUsername(username)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	// A person without credentials cannot log in.
	if person.Password == nil || *person.Password == "" {
		return nil, ErrInvalidCredentials
	}

	// Compare the plain text password against the stored
	// hash. The plain text is never stored.
	if bcrypt.CompareHashAndPassword(
		[]byte(*person.Password),
		[]byte(password),
	) != nil {
		return nil, ErrInvalidCredentials
	}

	return person, nil
}
