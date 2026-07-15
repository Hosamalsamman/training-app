package organizationtypes

type Service struct {
	repo Repository
}

func NewService() *Service {
	return &Service{
		repo: Repository{},
	}
}

func (s *Service) GetAll() ([]OrganizationType, error) {
	return s.repo.GetAll()
}

func (s *Service) GetByID(id int) (*OrganizationType, error) {
	return s.repo.GetByID(id)
}
