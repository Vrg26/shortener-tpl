package db

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

type dbPostgres struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *dbPostgres {
	return &dbPostgres{db: db}
}

func (p *dbPostgres) AddBatchURL(ctx context.Context, urls []ShortURL, userID uint32) ([]ShortURL, error) {
	tx, err := p.db.Begin()

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO urls (shorturl, originurl, userid) VALUES($1, $2, $3)")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	for index, url := range urls {
		id, err := p.generateID()

		if err != nil {
			return nil, err
		}
		if _, err = stmt.ExecContext(ctx, id, url.OriginURL, userID); err != nil {
			return nil, err
		}
		urls[index].ID = id
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func (p *dbPostgres) Add(ctx context.Context, url string, userID uint32) (string, error) {
	id, err := p.generateID()

	if err != nil {
		return "", err
	}

	_, err = p.db.ExecContext(ctx, "INSERT INTO urls (shorturl, originurl, userid) VALUES($1, $2, $3)", id, url, userID)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (p *dbPostgres) GetByOriginalURL(ctx context.Context, url string) (string, error) {
	row := p.db.QueryRowContext(ctx, "SELECT shorturl FROM urls WHERE originurl = $1", url)
	var result string
	if err := row.Scan(&result); err != nil {
		return "", err
	}
	return result, nil
}

func (p *dbPostgres) GetByURLAndUserID(ctx context.Context, url string, userID uint32) (ShortURL, error) {
	row := p.db.QueryRowContext(ctx, "SELECT shorturl, originurl, userid FROM urls WHERE originurl = $1 AND userid = $2", url, userID)
	var result ShortURL
	if err := row.Scan(&result.ID, &result.OriginURL, &result.UserID); err != nil {
		return result, err
	}
	return result, nil
}

func (p *dbPostgres) GetByID(ctx context.Context, id string) (ShortURL, error) {
	row := p.db.QueryRowContext(ctx, "SELECT shorturl, originurl FROM urls WHERE shorturl = $1", id)
	var result ShortURL
	if err := row.Scan(&result.ID, &result.OriginURL); err != nil {
		return result, err
	}
	return result, nil
}
func (p *dbPostgres) GetURLsByUserID(ctx context.Context, userID uint32) ([]ShortURL, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT shorturl, originurl, userid FROM urls WHERE userid = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []ShortURL

	for rows.Next() {
		var url ShortURL
		err = rows.Scan(&url.ID, &url.OriginURL, &url.UserID)
		if err != nil {
			return nil, err
		}
		result = append(result, url)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *dbPostgres) MigrateUp(sourceURL string) error {
	driver, err := postgres.WithInstance(p.db, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres", driver)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (p *dbPostgres) generateID() (string, error) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", buf), nil
}
