package db

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *sql.DB) error {
	return RunMigrationsWithPath(db, "file://migrations")
}

func RunMigrationsWithPath(db *sql.DB, migrationsPath string) error {
	migrationsDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrations, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", migrationsDriver)
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
