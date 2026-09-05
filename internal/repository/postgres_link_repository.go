package repository

import (
	"context"
	"database/sql"

	"github.com/SingletonVD/shortener/internal/model"
)

type PostgresLinkRepository struct {
	db *sql.DB
}

const (
	insertShortenedLinkQuery = "INSERT INTO shortened_links (short_link, full_link) VALUES($1, $2) ON CONFLICT (short_link) DO NOTHING"
)

func NewPostgresLinkRepository(db *sql.DB) (*PostgresLinkRepository, error) {
	return &PostgresLinkRepository{db: db}, nil
}

func (storage *PostgresLinkRepository) SaveIfAvailable(context context.Context, link model.ShortenedLink) (bool, error) {
	query := insertShortenedLinkQuery
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

func (storage *PostgresLinkRepository) SaveBatchIfAvailable(context context.Context, links []model.ShortenedLink) (bool, error) {
	tx, err := storage.db.BeginTx(context, nil)
	if err != nil {
		return false, err
	}

	statement, err := tx.PrepareContext(context, insertShortenedLinkQuery)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	for _, link := range links {
		result, err := statement.ExecContext(context, link.Short, link.FullLink)
		if err != nil {
			tx.Rollback()
			return false, err
		}

		insertedRows, err := result.RowsAffected()

		if err != nil {
			tx.Rollback()
			return false, err
		}

		if insertedRows != 1 {
			tx.Rollback()
			return false, nil
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return false, err
	}

	return true, nil
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
