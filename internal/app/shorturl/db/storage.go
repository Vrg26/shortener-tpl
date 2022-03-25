package db

type Storage interface {
	Add(url string, userID uint64) (string, error)
	GetByID(id string) (ShortURL, error)
	GetURLsByUserID(userId uint64) ([]ShortURL, error)
}
