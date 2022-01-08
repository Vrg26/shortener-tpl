package db

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/Vrg26/shortener-tpl/internal/app/shortUrl"
	"log"
	"sync"
)

var _ shortUrl.Storage = &db{}

type db struct {
	sync.Mutex
	urls map[string]shortUrl.ShortUrl
}

func (d *db) generateId() string {
	for {
		b := make([]byte, 8)
		_, err := rand.Read(b)
		if err != nil {
			log.Fatal(err)
		}
		newId := fmt.Sprintf("%x", b[0:8])
		if _, ok := d.urls[newId]; !ok {
			return newId
		}
	}
}

func (d *db) Add(url string) (string, error) {
	if d.urls == nil {
		d.Lock()
		d.urls = make(map[string]shortUrl.ShortUrl)
		d.Unlock()
	}

	newId := d.generateId()
	d.Lock()
	d.urls[newId] = shortUrl.ShortUrl{
		Id:        newId,
		OriginUrl: url,
	}
	d.Unlock()
	return newId, nil
}

func (d *db) GetById(id string) (shortUrl.ShortUrl, error) {
	if shortUrl, ok := d.urls[id]; ok {
		return shortUrl, nil
	}

	return shortUrl.ShortUrl{}, errors.New("short url not found")
}

func NewMemoryStorage() shortUrl.Storage {
	return &db{
		urls: make(map[string]shortUrl.ShortUrl),
	}
}
