package repository

import (
	"context"
	"database/sql"

	"github.com/SingletonVD/shortener/internal/model"
)

type PostgresLinkRepository struct {
	db *sql.DB
}

func NewPostgresLinkRepository(db *sql.DB) (*PostgresLinkRepository, error) {
	return &PostgresLinkRepository{db: db}, nil
}

func (storage *PostgresLinkRepository) SaveIfAvailable(context context.Context, link model.ShortenedLink) (bool, error) {
	query := "INSERT INTO shortened_links (short_link, full_link) VALUES($1, $2) ON CONFLICT (short_link) DO NOTHING"
	result, err := storage.db.ExecContext(context, query, link.Short, link.FullLink)
	if err != nil {
		return false, err
	}

	rowsInserted, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	return rowsInserted == 1, nil
}

func (storage *PostgresLinkRepository) FindLink(context context.Context, shortLink string) (*model.ShortenedLink, error) {
	query := "SELECT short_link, full_link FROM shortened_links WHERE short_link = $1 LIMIT 1"
	result := storage.db.QueryRowContext(context, query, shortLink)

	var shortenedLink model.ShortenedLink
	err := result.Scan(&shortenedLink.Short, &shortenedLink.FullLink)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &shortenedLink, nil
}
