package shorturl

import "github.com/Vrg26/shortener-tpl/internal/app/shorturl/db"

type Service struct {
	storage db.Storage
}

func NewService(st db.Storage) *Service {
	return &Service{
		storage: st,
	}
}

func (s *Service) Add(originURL string) (string, error) {
	newID, err := s.storage.Add(originURL)
	return newID, err
}

func (s *Service) GetByID(IDURL string) (db.ShortURL, error) {
	return s.storage.GetByID(IDURL)
}