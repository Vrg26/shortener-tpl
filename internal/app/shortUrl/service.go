package shortUrl

type Service struct {
	storage Storage
}

func NewService(st Storage) *Service {
	return &Service{
		storage: st,
	}
}

func (s *Service) Add(originUrl string) (string, error) {
	newId, err := s.storage.Add(originUrl)
	return newId, err
}

func (s *Service) GetById(IdUrl string) (ShortUrl, error) {
	return s.storage.GetById(IdUrl)
}
