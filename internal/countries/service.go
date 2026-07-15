package countries

type Service struct {
	repo Repository
}

func NewService() *Service {
	return &Service{
		repo: Repository{},
	}
}

func (s *Service) GetAll() ([]Country, error) {
	return s.repo.GetAll()
}

func (s *Service) GetByID(id int) (*Country, error) {
	return s.repo.GetByID(id)
}
