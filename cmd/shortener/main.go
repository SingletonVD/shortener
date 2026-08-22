package main

import (
	"log"
	"net/http"

	"github.com/SingletonVD/shortener/internal/config"
	"github.com/SingletonVD/shortener/internal/handler"
	"github.com/SingletonVD/shortener/internal/logger"
	"github.com/SingletonVD/shortener/internal/repository"
	"github.com/SingletonVD/shortener/internal/service"
)

func run() error {
	serverConfig := config.InitServerConfig()
	logger.InitializeLogger()
	linkStorage, err := repository.NewDiskLinkRepository(serverConfig.FileStoragePath)

	if err != nil {
		return err
	}

	linkService := service.NewLinkService(linkStorage)
	linkHandler := handler.NewLinkHandler(linkService, serverConfig.BaseLinkAddress)
	router := handler.NewRouter(linkHandler)

	return http.ListenAndServe(serverConfig.RunAddress, router)
}

func main() {
	err := run()
	if err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
	}
}
