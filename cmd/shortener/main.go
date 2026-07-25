package main

import (
	"net/http"

	"github.com/SingletonVD/shortener/internal/config"
	"github.com/SingletonVD/shortener/internal/handler"
	"github.com/SingletonVD/shortener/internal/repository"
	"github.com/SingletonVD/shortener/internal/service"
)

func run() error {
	serverConfig := config.InitServerConfig()
	linkStorage := repository.NewMemLinkRepository()
	linkService := service.NewLinkService(linkStorage)
	linkHandler := handler.NewLinkHandler(linkService, serverConfig.BaseLinkAddress)
	router := handler.NewRouter(linkHandler)

	return http.ListenAndServe(":8080", router)
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
