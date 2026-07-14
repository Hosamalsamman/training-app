package organizations

type Service struct {
	repo Repository
}

func NewService() *Service {
	return &Service{
		repo: Repository{},
	}
}

func (s *Service) GetAll() ([]Organization, error) {
	return s.repo.GetAll()
}

func (s *Service) GetByID(id string) (*Organization, error) {
	return s.repo.GetByID(id)
}
