package db

type Storage interface {
	Add(url string) (string, error)
	GetById(id string) (ShortUrl, error)
}
