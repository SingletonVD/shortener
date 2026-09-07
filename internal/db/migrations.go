package db

import (
	"database/sql"
	"errors"

	"github.com/SingletonVD/shortener/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func RunMigrations(db *sql.DB) error {
	src, err := iofs.New(migrations.MigrationsFS, ".")
	if err != nil {
		return err
	}

	migrationsDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrations, err := migrate.NewWithInstance("iofs", src, "postgres", migrationsDriver)
	if err != nil {
		return err
	}

	err = migrations.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}
	return nil
}
