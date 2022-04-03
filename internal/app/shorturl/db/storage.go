package db

import (
	"context"
)

type Storage interface {
	Add(ctx context.Context, url string, userID uint32) (string, error)
	GetByID(ctx context.Context, id string) (ShortURL, error)
	GetByOriginalURL(ctx context.Context, url string) (string, error)
	GetURLsByUserID(ctx context.Context, userID uint32) ([]ShortURL, error)
	AddBatchURL(ctx context.Context, urls []ShortURL, userID uint32) ([]ShortURL, error)
	DeleteURLs(ctx context.Context, ids []string, userID uint32) error
}
