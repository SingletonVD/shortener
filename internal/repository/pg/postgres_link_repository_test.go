package pg

import (
	"testing"

	"github.com/SingletonVD/shortener/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func TestPostgresLinkRepository(t *testing.T) {

	ctx := t.Context()
	setupTest := func() {
		_, err := dbPool.Exec(ctx, "DELETE FROM shortened_links")
		require.NoError(t, err)
	}

	postgresLinkRepository, err := NewPostgresLinkRepository(dbPool)
	require.NoError(t, err)

	t.Run("SaveIfAvailable", func(t *testing.T) {
		setupTest()

		shortenedLink := model.ShortenedLink{
			Short:    "12345678",
			FullLink: "https://yandex.ru",
		}

		saved, err := postgresLinkRepository.SaveIfAvailable(ctx, shortenedLink)
		require.NoError(t, err)
		require.Equal(t, true, saved)

		_, err = postgresLinkRepository.SaveIfAvailable(ctx, shortenedLink)
		require.EqualError(t, err, (&FullLinkConflict{ShortLink: "12345678", FullLink: "https://yandex.ru"}).Error())

		shortenedLink = model.ShortenedLink{
			Short:    "12345678",
			FullLink: "https://yandex2.ru",
		}

		saved, err = postgresLinkRepository.SaveIfAvailable(ctx, shortenedLink)
		require.NoError(t, err)
		require.Equal(t, false, saved)
	})

	t.Run("FindLink", func(t *testing.T) {
		setupTest()

		shortenedLink := model.ShortenedLink{
			Short:    "01234567",
			FullLink: "https://yandex2.ru",
		}

		_, err := postgresLinkRepository.SaveIfAvailable(ctx, shortenedLink)
		require.NoError(t, err)

		foundLink, err := postgresLinkRepository.FindLink(ctx, shortenedLink.Short)
		require.NoError(t, err)
		require.NotNil(t, foundLink)
		require.Equal(t, shortenedLink, *foundLink)

		notFoundLink, err := postgresLinkRepository.FindLink(ctx, "non_existent")
		require.NoError(t, err)
		require.Nil(t, notFoundLink)

	})

	t.Run("SaveBatchIfAvailable", func(t *testing.T) {
		setupTest()

		shortenedLinks := []model.ShortenedLink{{
			Short:    "23456789",
			FullLink: "https://yandex3.ru",
		}, {
			Short:    "34567890",
			FullLink: "https://ya.ru",
		}}

		saved, err := postgresLinkRepository.SaveBatchIfAvailable(ctx, shortenedLinks)
		require.NoError(t, err)
		require.Equal(t, true, saved)

		foundLink, err := postgresLinkRepository.FindLink(ctx, "23456789")
		require.NoError(t, err)
		require.NotNil(t, foundLink)
		require.Equal(t, "https://yandex3.ru", foundLink.FullLink)

		notSavedShortenedLinks := []model.ShortenedLink{{
			Short:    "23456789",
			FullLink: "https://yandex4.ru",
		}, {
			Short:    "45678901",
			FullLink: "https://ya2.ru",
		}}

		saved, err = postgresLinkRepository.SaveBatchIfAvailable(ctx, notSavedShortenedLinks)
		require.NoError(t, err)
		require.Equal(t, false, saved)

		notFoundLink, err := postgresLinkRepository.FindLink(ctx, "45678901")
		require.NoError(t, err)
		require.Nil(t, notFoundLink)
	})
}
