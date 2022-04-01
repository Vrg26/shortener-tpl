package db

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"os"
)

type dbFile struct {
	filePath string
}

func NewFileStorage(filePath string) *dbFile {
	return &dbFile{
		filePath: filePath,
	}
}

func (f *dbFile) AddBatchURL(ctx context.Context, urls []ShortURL, userID uint32) ([]ShortURL, error) {
	for index, url := range urls {
		id, err := f.Add(ctx, url.OriginURL, userID)
		if err != nil {
			return nil, err
		}
		urls[index].ID = id
	}
	return urls, nil
}

func (f *dbFile) Add(ctx context.Context, url string, userID uint32) (string, error) {

	shortURL, err := f.GetByURLAndUserID(url, userID)
	if err != nil {
		return "", err
	}

	if shortURL.ID != "" {
		return shortURL.ID, nil
	}

	newID := f.generateID()

	sURL := ShortURL{
		ID:        newID,
		OriginURL: url,
		UserID:    userID,
	}

	file, err := os.OpenFile(f.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := json.Marshal(&sURL)
	if err != nil {
		return "", err
	}

	writer := bufio.NewWriter(file)
	if _, err := writer.Write(data); err != nil {
		return "", err
	}
	if err := writer.WriteByte('\n'); err != nil {
		return "", err
	}
	if err := writer.Flush(); err != nil {
		return "", err
	}

	return newID, nil
}

func (f *dbFile) GetByOriginalURL(ctx context.Context, url string) (string, error) {
	file, err := os.OpenFile(f.filePath, os.O_RDONLY|os.O_CREATE, 0777)

	if err != nil {
		return "", err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var sURL ShortURL
		err := json.Unmarshal(scanner.Bytes(), &sURL)
		if err != nil {
			return "", err
		}
		if sURL.OriginURL == url {
			return sURL.ID, nil
		}
	}
	return "", nil
}

func (f *dbFile) GetURLsByUserID(ctx context.Context, userID uint32) ([]ShortURL, error) {
	file, err := os.OpenFile(f.filePath, os.O_RDONLY, 0777)

	if err != nil {
		return nil, err
	}
	defer file.Close()
	var resultUrls []ShortURL
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var sURL ShortURL
		err := json.Unmarshal(scanner.Bytes(), &sURL)
		if err != nil {
			return nil, err
		}
		if sURL.UserID == userID {
			resultUrls = append(resultUrls, sURL)
		}
	}

	return resultUrls, nil
}

func (f *dbFile) GetByURLAndUserID(url string, userID uint32) (ShortURL, error) {
	file, err := os.OpenFile(f.filePath, os.O_RDONLY|os.O_CREATE, 0777)

	if err != nil {
		return ShortURL{}, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var sURL ShortURL
		err := json.Unmarshal(scanner.Bytes(), &sURL)
		if err != nil {
			return ShortURL{}, err
		}
		if sURL.OriginURL == url && sURL.UserID == userID {
			return sURL, nil
		}
	}
	return ShortURL{}, nil
}

func (f *dbFile) GetByID(ctx context.Context, id string) (ShortURL, error) {
	file, err := os.OpenFile(f.filePath, os.O_RDONLY|os.O_CREATE, 0777)

	if err != nil {
		return ShortURL{}, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var sURL ShortURL
		err := json.Unmarshal(scanner.Bytes(), &sURL)
		if err != nil {
			return ShortURL{}, err
		}
		if sURL.ID == id {
			return sURL, nil
		}
	}

	return ShortURL{}, errors.New("short url not found")
}

func (f *dbFile) generateID() string {
	lenID := 6
	chars := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_")
	uniqueRuneArray := make([]rune, lenID)
	for i := range uniqueRuneArray {
		uniqueRuneArray[i] = chars[rand.Intn(len(chars))]
	}
	return string(uniqueRuneArray)
}
