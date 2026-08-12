package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	RunAddress      string
	BaseLinkAddress string
}

func InitServerConfig() *ServerConfig {
	serverConfig := new(ServerConfig)
	flag.StringVar(&serverConfig.RunAddress, "a", "localhost:8080", "address to bind server")
	flag.StringVar(&serverConfig.BaseLinkAddress, "b", "http://localhost:8080", "base address for shortened link")
	flag.Parse()

	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		serverConfig.RunAddress = envServerAddress
	}

	if envBaseUrl := os.Getenv("BASE_URL"); envBaseUrl != "" {
		serverConfig.BaseLinkAddress = envBaseUrl
	}

	return serverConfig
}
