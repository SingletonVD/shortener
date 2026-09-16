package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	intDB "github.com/SingletonVD/shortener/internal/db"
	"github.com/SingletonVD/shortener/internal/model"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type PostgresLinkRepository struct {
	pool *pgxpool.Pool
}

type FullLinkConflict struct {
	ShortLink string
	FullLink  string
}

func (conflict *FullLinkConflict) Error() string {
	return fmt.Sprintf("Full link %s already exists as %s", conflict.FullLink, conflict.ShortLink)
}

func NewPostgresLinkRepository(pool *pgxpool.Pool) (*PostgresLinkRepository, error) {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()
	err := intDB.RunMigrations(db)
	if err != nil {
		return nil, err
	}
	return &PostgresLinkRepository{pool: pool}, nil
}

func (storage *PostgresLinkRepository) SaveIfAvailable(ctx context.Context, link model.ShortenedLink) (bool, error) {
	query := "INSERT INTO shortened_links (short_link, full_link) VALUES($1, $2) ON CONFLICT (full_link) DO NOTHING"
	tx, err := storage.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, query, link.Short, link.FullLink)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				if pgErr.ConstraintName == "idx_shortened_links_short_link" {
					return false, nil
				}
			}
		}
		return false, err
	}

	rowsInserted := result.RowsAffected()
	if rowsInserted != 1 {
		query := "SELECT short_link, full_link FROM shortened_links WHERE full_link = $1 LIMIT 1"
		result := tx.QueryRow(ctx, query, link.FullLink)

		var conflict FullLinkConflict
		err := result.Scan(&conflict.ShortLink, &conflict.FullLink)
		if err != nil {
			return false, err
		}
		return false, &conflict
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (storage *PostgresLinkRepository) SaveBatchIfAvailable(ctx context.Context, links []model.ShortenedLink) (bool, error) {
	query := "INSERT INTO shortened_links (short_link, full_link) VALUES($1, $2) ON CONFLICT (short_link) DO NOTHING"
	tx, err := storage.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}

	for _, link := range links {
		batch.Queue(query, link.Short, link.FullLink)
	}

	batchResults := tx.SendBatch(ctx, batch)
	defer batchResults.Close()

	for range len(links) {
		result, err := batchResults.Exec()
		if err != nil {
			return false, err
		}

		if result.RowsAffected() != 1 {
			return false, nil
		}
	}

	err = batchResults.Close()
	if err != nil {
		return false, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (storage *PostgresLinkRepository) FindLink(ctx context.Context, shortLink string) (*model.ShortenedLink, error) {
	query := "SELECT short_link, full_link FROM shortened_links WHERE short_link = $1 LIMIT 1"
	result := storage.pool.QueryRow(ctx, query, shortLink)

	var shortenedLink model.ShortenedLink
	err := result.Scan(&shortenedLink.Short, &shortenedLink.FullLink)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &shortenedLink, nil
}

func (storage *PostgresLinkRepository) Ping(ctx context.Context) error {
	return storage.pool.Ping(ctx)
}
