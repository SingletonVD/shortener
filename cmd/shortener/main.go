package main

import (
	"log"
	"net/http"

	"github.com/SingletonVD/shortener/internal/config"
	"github.com/SingletonVD/shortener/internal/handler"
	"github.com/SingletonVD/shortener/internal/logger"
	"github.com/SingletonVD/shortener/internal/repository"
	"github.com/SingletonVD/shortener/internal/service"

	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func run() error {
	serverConfig := config.InitServerConfig()
	logger.InitializeLogger()

	db, err := sql.Open("pgx", serverConfig.DatabaseDsn)
	if err != nil {
		return err
	}
	defer db.Close()

	linkStorage, err := repository.NewLinkRepository(serverConfig, db)

	if err != nil {
		return err
	}

	linkService := service.NewLinkService(linkStorage)
	linkHandler := handler.NewLinkHandler(linkService, serverConfig.BaseLinkAddress)
	pingHandler := handler.NewPingHandler(db)
	router := handler.NewRouter(linkHandler, pingHandler)

	return http.ListenAndServe(serverConfig.RunAddress, router)
}

func main() {
	err := run()
	if err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
	}
}
