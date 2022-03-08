package db

import (
	"bufio"
	"encoding/json"
	"errors"
	"math/rand"
	"os"
)

var _ Storage = &dbFile{}

type dbFile struct {
	filePath string
}

func (f *dbFile) Add(url string) (string, error) {

	newID, _ := f.GetByURL(url)
	if newID != "" {
		return newID, nil
	}
	newID = f.generateID()

	sURL := ShortURL{
		ID:        newID,
		OriginURL: url,
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

func (f *dbFile) GetByURL(url string) (string, error) {
	file, err := os.OpenFile(f.filePath, os.O_RDONLY|os.O_CREATE, 0777)

	if err != nil {
		return "", err
	}

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

func (f *dbFile) GetByID(id string) (ShortURL, error) {
	file, err := os.OpenFile(f.filePath, os.O_RDONLY|os.O_CREATE, 0777)

	if err != nil {
		return ShortURL{}, err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var shortUrl ShortURL
		err := json.Unmarshal(scanner.Bytes(), &shortUrl)
		if err != nil {
			return ShortURL{}, err
		}
		if shortUrl.ID == id {
			return shortUrl, nil
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

func NewFileStorage(filePath string) (Storage, error) {
	return &dbFile{
		filePath: filePath,
	}, nil
}
