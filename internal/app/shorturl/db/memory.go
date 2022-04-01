package db

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"sync"
)

type dbMemory struct {
	sync.Mutex
	urls map[string]ShortURL
}

func NewMemoryStorage() *dbMemory {
	return &dbMemory{
		urls: make(map[string]ShortURL),
	}
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

func (d *dbMemory) GetByOriginalURL(ctx context.Context, url string) (string, error) {
	for _, itemMap := range d.urls {
		if itemMap.OriginURL == url {
			return itemMap.ID, nil
		}
	}
	return "", errors.New("short url not found")
}

func (d *dbMemory) GetURLsByUserID(ctx context.Context, userID uint32) ([]ShortURL, error) {
	if d.urls == nil {
		return []ShortURL{}, nil
	}
	var resultURLs []ShortURL
	for _, itemMap := range d.urls {
		if itemMap.UserID == userID {
			resultURLs = append(resultURLs, itemMap)
		}
	}
	return resultURLs, nil
}

func (d *dbMemory) AddBatchURL(ctx context.Context, urls []ShortURL, userID uint32) ([]ShortURL, error) {
	for index, url := range urls {
		id, err := d.Add(ctx, url.OriginURL, userID)
		if err != nil {
			return nil, err
		}
		urls[index].ID = id
	}
	return urls, nil
}

func (d *dbMemory) Add(ctx context.Context, url string, userID uint32) (string, error) {
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
		UserID:    userID,
	}
	d.Unlock()
	return newID, nil
}

func (d *dbMemory) GetByID(ctx context.Context, id string) (ShortURL, error) {
	if ShortURL, ok := d.urls[id]; ok {
		return ShortURL, nil
	}

	return ShortURL{}, errors.New("short url not found")
}
