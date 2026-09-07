package pg

import (
	"context"
	"log"
	"os"
	"testing"

	dbInternal "github.com/SingletonVD/shortener/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var dbPool *pgxpool.Pool

const (
	pgImageName = "postgres:18-alpine"
	dbName      = "shortener"
	user        = "user"
	password    = "password"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	ctr, err := postgres.Run(
		ctx,
		pgImageName,
		postgres.WithDatabase(dbName),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	if err != nil {
		log.Fatalf("Failed to start container string")
	}
	defer ctr.Terminate(ctx)

	dbURL, err := ctr.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("Failed to get connection string")
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to create connection pool")
	}
	dbPool = pool

	db := stdlib.OpenDBFromPool(pool)
	err = dbInternal.RunMigrations(db)
	if err != nil {
		log.Fatalf("Failed to run migrations connection string")
	}

	code := m.Run()
	os.Exit(code)
}
