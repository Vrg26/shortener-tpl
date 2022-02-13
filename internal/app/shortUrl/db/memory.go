package db

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"sync"
)

var _ Storage = &db{}

type db struct {
	sync.Mutex
	urls map[string]ShortUrl
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
		d.urls = make(map[string]ShortUrl)
		d.Unlock()
	}

	newId := d.generateId()
	d.Lock()
	d.urls[newId] = ShortUrl{
		Id:        newId,
		OriginUrl: url,
	}
	d.Unlock()
	return newId, nil
}

func (d *db) GetById(id string) (ShortUrl, error) {
	if shortUrl, ok := d.urls[id]; ok {
		return shortUrl, nil
	}

	return ShortUrl{}, errors.New("short url not found")
}

func NewMemoryStorage() Storage {
	return &db{
		urls: make(map[string]ShortUrl),
	}
}
