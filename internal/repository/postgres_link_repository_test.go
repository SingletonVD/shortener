package repository

import (
	"context"
	"database/sql"
	"testing"

	dbInternal "github.com/SingletonVD/shortener/internal/db"
	"github.com/SingletonVD/shortener/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestPostgresLinkRepository(t *testing.T) {

	// это стоит унести и сделать какой-нибудь синглтон, но пока контейнер нужен только для тестов одного репозитория, имхо будет только усложнением
	context := context.Background()

	dbName := "shortener"
	user := "user"
	password := "password"

	ctr, err := postgres.Run(
		context,
		"postgres:18-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	dbURL, err := ctr.ConnectionString(context)
	require.NoError(t, err)

	db, err := sql.Open("pgx", dbURL)
	require.NoError(t, err)

	err = dbInternal.RunMigrationsWithPath(db, "file://../../migrations")
	require.NoError(t, err)

	setupTest := func() {
		t.Cleanup(func() {
			_, err := db.ExecContext(context, "DELETE FROM shortened_links")
			require.NoError(t, err)
		})
	}

	postgresLinkRepository, err := NewPostgresLinkRepository(db)
	require.NoError(t, err)

	t.Run("SaveIfAvailable", func(t *testing.T) {
		setupTest()

		shortenedLink := model.ShortenedLink{
			Short:    "12345678",
			FullLink: "https://yandex.ru",
		}

		saved, err := postgresLinkRepository.SaveIfAvailable(context, shortenedLink)
		require.NoError(t, err)
		require.Equal(t, true, saved)

		saved, err = postgresLinkRepository.SaveIfAvailable(context, shortenedLink)
		require.NoError(t, err)
		require.Equal(t, false, saved)
	})

	t.Run("FindLink", func(t *testing.T) {
		setupTest()

		shortenedLink := model.ShortenedLink{
			Short:    "12345678",
			FullLink: "https://yandex.ru",
		}

		_, err := postgresLinkRepository.SaveIfAvailable(context, shortenedLink)
		require.NoError(t, err)

		foundLink, err := postgresLinkRepository.FindLink(context, shortenedLink.Short)
		require.NoError(t, err)
		require.NotNil(t, foundLink)
		require.Equal(t, shortenedLink, *foundLink)

		notFoundLink, err := postgresLinkRepository.FindLink(context, "non_existent")
		require.NoError(t, err)
		require.Nil(t, notFoundLink)

	})

	t.Run("SaveBatchIfAvailable", func(t *testing.T) {
		setupTest()

		shortenedLinks := []model.ShortenedLink{{
			Short:    "23456789",
			FullLink: "https://yandex.ru",
		}, {
			Short:    "34567890",
			FullLink: "https://ya.ru",
		}}

		saved, err := postgresLinkRepository.SaveBatchIfAvailable(context, shortenedLinks)
		require.NoError(t, err)
		require.Equal(t, true, saved)

		foundLink, err := postgresLinkRepository.FindLink(context, "23456789")
		require.NoError(t, err)
		require.NotNil(t, foundLink)
		require.Equal(t, "https://yandex.ru", foundLink.FullLink)

		notSavedShortenedLinks := []model.ShortenedLink{{
			Short:    "23456789",
			FullLink: "https://yandex.ru",
		}, {
			Short:    "45678901",
			FullLink: "https://ya.ru",
		}}

		saved, err = postgresLinkRepository.SaveBatchIfAvailable(context, notSavedShortenedLinks)
		require.NoError(t, err)
		require.Equal(t, false, saved)

		notFoundLink, err := postgresLinkRepository.FindLink(context, "45678901")
		require.NoError(t, err)
		require.Nil(t, notFoundLink)
	})
}
