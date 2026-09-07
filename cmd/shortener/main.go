package main

import (
	"context"
	"log"
	"net/http"

	"github.com/SingletonVD/shortener/internal/config"
	"github.com/SingletonVD/shortener/internal/handler"
	"github.com/SingletonVD/shortener/internal/logger"
	"github.com/SingletonVD/shortener/internal/repository/disk"
	"github.com/SingletonVD/shortener/internal/repository/memory"
	"github.com/SingletonVD/shortener/internal/repository/pg"
	"github.com/SingletonVD/shortener/internal/service"

	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type AppContainer struct {
	pinger         handler.Pinger
	linkRepository service.LinkRepository
}

func newAppContainer(ctx context.Context, serverConfig *config.ServerConfig) (*AppContainer, error) {
	if serverConfig.DatabaseDsn != "" {
		pool, err := pgxpool.New(ctx, serverConfig.DatabaseDsn)
		if err != nil {
			return nil, err
		}
		repository, err := pg.NewPostgresLinkRepository(pool)
		if err != nil {
			return nil, err
		}
		return &AppContainer{
			pinger:         repository,
			linkRepository: repository,
		}, nil
	}

	if serverConfig.FileStoragePath != "" {
		repository, err := disk.NewDiskLinkRepository(serverConfig.FileStoragePath)
		if err != nil {
			return nil, err
		}
		return &AppContainer{
			pinger:         repository,
			linkRepository: repository,
		}, nil
	}

	repository := memory.NewMemLinkRepository()
	return &AppContainer{
		pinger:         repository,
		linkRepository: repository,
	}, nil
}

func run() error {
	serverConfig := config.InitServerConfig()
	logger.InitializeLogger()

	db, err := sql.Open("pgx", serverConfig.DatabaseDsn)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx := context.Background()
	appContainer, err := newAppContainer(ctx, serverConfig)

	if err != nil {
		return err
	}

	linkService := service.NewLinkService(appContainer.linkRepository)
	linkHandler := handler.NewLinkHandler(linkService, serverConfig.BaseLinkAddress)
	pingHandler := handler.NewPingHandler(appContainer.pinger)
	router := handler.NewRouter(linkHandler, pingHandler)

	return http.ListenAndServe(serverConfig.RunAddress, router)
}

func main() {
	err := run()
	if err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
	}
}
