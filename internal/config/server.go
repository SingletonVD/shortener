package config

import "flag"

type ServerConfig struct {
	RunAddress      string
	BaseLinkAddress string
}

func InitServerConfig() *ServerConfig {
	serverConfig := new(ServerConfig)
	flag.StringVar(&serverConfig.RunAddress, "a", "localhost:8080", "address to bind server")
	flag.StringVar(&serverConfig.BaseLinkAddress, "b", "http://localhost:8080", "base address for shortened link")
	flag.Parse()
	return serverConfig
}
