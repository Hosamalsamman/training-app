package organizationtypes

type Service struct {
	repo Repository
}

func NewService() *Service {
	return &Service{
		repo: Repository{},
	}
}

func (s *Service) GetOrganizationTypes() ([]OrganizationType, error) {
	return s.repo.GetAllOrganizationTypes()
}

func (s *Service) GetOrganizationType(id string) (*OrganizationType, error) {
	return s.repo.GetOrganizationTypeByID(id)
}
