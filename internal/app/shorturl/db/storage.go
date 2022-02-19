package db

type Storage interface {
	Add(url string) (string, error)
	GetByID(id string) (ShortURL, error)
}
