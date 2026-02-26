package config

import (
	"os"
)

type Config struct {
	ServerPort string
}

func Load() (*Config, error) {
	serverPort := os.Getenv("SERVER_PORT")

	return &Config{
		ServerPort: serverPort,
	}, nil
}
