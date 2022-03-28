package shorturl

import (
	"context"
	"github.com/Vrg26/shortener-tpl/internal/app/shorturl/db"
)

type Service struct {
	storage db.Storage
}

func NewService(st db.Storage) *Service {
	return &Service{
		storage: st,
	}
}

func (s *Service) Add(ctx context.Context, originURL string, userId uint32) (string, error) {
	newID, err := s.storage.Add(ctx, originURL, userId)
	return newID, err
}

func (s *Service) GetURLsByUserID(ctx context.Context, userId uint32) ([]db.ShortURL, error) {
	return s.storage.GetURLsByUserID(ctx, userId)
}

func (s *Service) GetByID(ctx context.Context, IDURL string) (db.ShortURL, error) {
	return s.storage.GetByID(ctx, IDURL)
}
