package db

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"sync"
)

var _ Storage = &dbMemory{}

type dbMemory struct {
	sync.Mutex
	urls map[string]ShortURL
}

func (d *dbMemory) generateID() string {
	for {
		b := make([]byte, 8)
		_, err := rand.Read(b)
		if err != nil {
			log.Fatal(err)
		}
		newID := fmt.Sprintf("%x", b[0:8])
		if _, ok := d.urls[newID]; !ok {
			return newID
		}
	}
}

func (d *dbMemory) Add(url string) (string, error) {
	if d.urls == nil {
		d.Lock()
		d.urls = make(map[string]ShortURL)
		d.Unlock()
	}

	newID := d.generateID()
	d.Lock()
	d.urls[newID] = ShortURL{
		ID:        newID,
		OriginURL: url,
	}
	d.Unlock()
	return newID, nil
}

func (d *dbMemory) GetByID(id string) (ShortURL, error) {
	if ShortURL, ok := d.urls[id]; ok {
		return ShortURL, nil
	}

	return ShortURL{}, errors.New("short url not found")
}

func NewMemoryStorage() Storage {
	return &dbMemory{
		urls: make(map[string]ShortURL),
	}
}
