package config

import (
	"flag"
	"os"

	"github.com/SingletonVD/shortener/internal/worker"
)

type ServerConfig struct {
	RunAddress         string
	BaseLinkAddress    string
	FileStoragePath    string
	DatabaseDsn        string
	AuthSecret         string
	DeleteWorkerConfig worker.DeleteWorkerConfig
}

func InitServerConfig() *ServerConfig {
	serverConfig := new(ServerConfig)
	flag.StringVar(&serverConfig.RunAddress, "a", "localhost:8080", "address to bind server")
	flag.StringVar(&serverConfig.BaseLinkAddress, "b", "http://localhost:8080", "base address for shortened link")
	flag.StringVar(&serverConfig.FileStoragePath, "f", "", "path to json storage file")
	flag.StringVar(&serverConfig.DatabaseDsn, "d", "", "database connection address")
	flag.StringVar(&serverConfig.AuthSecret, "s", "", "jwt auth secret")
	flag.IntVar(&serverConfig.DeleteWorkerConfig.BatchSize, "dwBatchSize", 100, "delete worker batch size")
	flag.IntVar(&serverConfig.DeleteWorkerConfig.BufferSize, "dwBufferSize", 1024, "delete worker buffer size")
	flag.IntVar(&serverConfig.DeleteWorkerConfig.ConcurrentWriters, "dwConcurrentWriters", 10, "delete worker concurrent queue writers")
	flag.IntVar(&serverConfig.DeleteWorkerConfig.ScheduleInterval, "dwScheduleInterval", 10, "delete worker schedule interval in seconds")
	flag.Parse()

	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		serverConfig.RunAddress = envServerAddress
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		serverConfig.BaseLinkAddress = envBaseURL
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		serverConfig.FileStoragePath = envFileStoragePath
	}

	if databaseDsn := os.Getenv("DATABASE_DSN"); databaseDsn != "" {
		serverConfig.DatabaseDsn = databaseDsn
	}

	if authSecret := os.Getenv("AUTH_SECRET"); authSecret != "" {
		serverConfig.AuthSecret = authSecret
	}

	return serverConfig
}
