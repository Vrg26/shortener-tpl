package shortUrl

import "github.com/Vrg26/shortener-tpl/internal/app/shortUrl/db"

type Service struct {
	storage db.Storage
}

func NewService(st db.Storage) *Service {
	return &Service{
		storage: st,
	}
}

func (s *Service) Add(originUrl string) (string, error) {
	newId, err := s.storage.Add(originUrl)
	return newId, err
}

func (s *Service) GetById(IdUrl string) (db.ShortUrl, error) {
	return s.storage.GetById(IdUrl)
}
