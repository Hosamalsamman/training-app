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

// GetAllUsers returns the client's registered users: the
// persons who already have credentials. The frontend uses it
// to render the users an admin may pick as a reset-password
// target, so he cannot select a person without a user.
func (s *Service) GetAllUsers(clientID int) ([]models.UserSummary, error) {

	return s.repo.ForClient(clientID).GetAllUsers()
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

// AdminGroupID is the group id assigned to administrator
// users. Only admins may reset other users' passwords.
const AdminGroupID = 1

// Sentinel errors for the register-user, password flows.
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
	// The username became empty after normalization.
	ErrInvalidUsername = errors.New("username is empty")
	// Login failed: wrong username or wrong password.
	// Deliberately vague so the response cannot reveal
	// which part was wrong.
	ErrInvalidCredentials = errors.New("invalid username or password")
	// The supplied old password did not match the stored hash.
	ErrWrongOldPassword = errors.New("old password is incorrect")
	// The person has no credentials yet, so there is no
	// password to change or reset.
	ErrNotRegistered = errors.New("person is not registered as a user")
	// The caller is not an admin, so he may not reset
	// other users' passwords.
	ErrNotAdmin = errors.New("only admins can reset passwords")
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

// ChangePassword lets an authenticated user change his own
// password. The old password must match the stored hash, so
// only someone who knows the current one can perform the
// change. No transaction is needed: the flow is one read and
// one independent single-column update.
func (s *Service) ChangePassword(
	clientID int,
	userID int,
	req *models.ChangePasswordRequest,
) error {

	// Restrict every lookup and the update to the
	// authenticated client's data.
	repo := s.repo.ForClient(clientID)

	// The current hash is needed to verify the old password.
	current, err := repo.GetPasswordByID(userID)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}

		return err
	}

	// A nil password means the person was never registered
	// as a user, so there is nothing to change here.
	if current == nil {
		return ErrNotRegistered
	}

	// The old password must match, otherwise anyone holding
	// a valid token could take over the account.
	if bcrypt.CompareHashAndPassword(
		[]byte(*current),
		[]byte(req.OldPassword),
	) != nil {
		return ErrWrongOldPassword
	}

	// Hash the new password. The plain text is never stored.
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	return repo.UpdatePassword(userID, string(hashedPassword))
}

// ResetPassword lets an admin (group 1) set a new password
// for another user without knowing the old one. This is the
// recovery path for users who forgot their password. The
// caller's group is read from the database, not from the
// token, so a demoted admin cannot reuse an old token.
func (s *Service) ResetPassword(
	clientID int,
	callerID int,
	targetPersonID int,
	req *models.ResetPasswordRequest,
) error {

	// Restrict every lookup and the update to the
	// authenticated client's data.
	repo := s.repo.ForClient(clientID)

	// Only admins may reset passwords.
	callerGroupID, err := repo.GetGroupID(callerID)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}

		return err
	}

	if callerGroupID == nil || *callerGroupID != AdminGroupID {
		return ErrNotAdmin
	}

	// The target must exist in this client and already be a
	// registered user; RegisterUser is the flow that grants
	// initial credentials.
	target, err := repo.GetPasswordByID(targetPersonID)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPersonNotFound
		}

		return err
	}

	if target == nil {
		return ErrNotRegistered
	}

	// Hash the new password. The plain text is never stored.
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	return repo.UpdatePassword(targetPersonID, string(hashedPassword))
}
