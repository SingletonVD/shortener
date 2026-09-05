package repository

import (
	"context"
	"database/sql"

	"github.com/SingletonVD/shortener/internal/config"
	intDB "github.com/SingletonVD/shortener/internal/db"
	"github.com/SingletonVD/shortener/internal/model"
)

type LinkRepository interface {
	SaveIfAvailable(context context.Context, link model.ShortenedLink) (bool, error)
	FindLink(context context.Context, shortLink string) (*model.ShortenedLink, error)
}

func NewLinkRepository(serverConfig *config.ServerConfig, db *sql.DB) (LinkRepository, error) {
	if serverConfig.DatabaseDsn != "" {
		/* это явно стоит делать не при инициализации репозитория, но пока что не нашел места лучше
		когда/если будет несколько репозиториев, скорее всего вынесу в какую-нибудь обертку повыше
		*/
		err := intDB.RunMigrations(db)

		if err != nil {
			return nil, err
		}

		return NewPostgresLinkRepository(db)
	}

	if serverConfig.FileStoragePath != "" {
		return NewDiskLinkRepository(serverConfig.FileStoragePath)
	}

	return NewMemLinkRepository(), nil
}
