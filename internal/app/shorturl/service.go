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

func (s *Service) Add(ctx context.Context, originURL string, userID uint32) (string, error) {
	newID, err := s.storage.Add(ctx, originURL, userID)
	return newID, err
}

func (s *Service) AddBatchURL(ctx context.Context, urls []db.ShortURL, userID uint32) ([]db.ShortURL, error) {
	return s.storage.AddBatchURL(ctx, urls, userID)
}

func (s *Service) GetURLsByUserID(ctx context.Context, userID uint32) ([]db.ShortURL, error) {
	return s.storage.GetURLsByUserID(ctx, userID)
}

func (s *Service) GetByOriginalURL(ctx context.Context, url string) (string, error) {
	return s.storage.GetByOriginalURL(ctx, url)
}

func (s *Service) GetByID(ctx context.Context, idURL string) (db.ShortURL, error) {
	return s.storage.GetByID(ctx, idURL)
}
