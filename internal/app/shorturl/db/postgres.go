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

func (p *dbPostgres) Add(ctx context.Context, url string, userId uint32) (string, error) {
	id, err := p.generateID()

	if err != nil {
		return "", err
	}
	shorturl, err := p.GetByURLAndUserId(ctx, url, userId)

	if err == nil {
		return shorturl.ID, nil
	}

	_, err = p.db.ExecContext(ctx, "INSERT INTO urls (shorturl, originurl, userid) VALUES($1, $2, $3)", id, url, userId)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (p *dbPostgres) GetByURLAndUserId(ctx context.Context, url string, userId uint32) (ShortURL, error) {
	row := p.db.QueryRowContext(ctx, "SELECT shorturl, originurl, userid FROM urls WHERE originurl = $1 AND userid = $2", url, userId)
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
func (p *dbPostgres) GetURLsByUserID(ctx context.Context, userId uint32) ([]ShortURL, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT shorturl, originurl, userid FROM urls WHERE userid = $1", userId)
	var result []ShortURL
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var url ShortURL
		err = rows.Scan(&url.ID, &url.OriginURL, &url.UserID)
		if err != nil {
			return nil, err
		}
		result = append(result, url)
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
