package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	RunAddress      string
	BaseLinkAddress string
	FileStoragePath string
}

func InitServerConfig() *ServerConfig {
	serverConfig := new(ServerConfig)
	flag.StringVar(&serverConfig.RunAddress, "a", "localhost:8080", "address to bind server")
	flag.StringVar(&serverConfig.BaseLinkAddress, "b", "http://localhost:8080", "base address for shortened link")
	flag.StringVar(&serverConfig.FileStoragePath, "f", "db.json", "path to json storage file")
	flag.Parse()

	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		serverConfig.RunAddress = envServerAddress
	}

	if envBaseUrl := os.Getenv("BASE_URL"); envBaseUrl != "" {
		serverConfig.BaseLinkAddress = envBaseUrl
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		serverConfig.FileStoragePath = envFileStoragePath
	}

	return serverConfig
}
