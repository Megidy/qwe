package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	// it is not used, but should have been
	LogLevel string `env:"LOG_LEVEL"`

	HttpServerPort string `env:"HTTP_SERVER_PORT,required"`
	PostgresURI    string `env:"POSTGRES_URI,required" `
}

func NewConfig() (*Config, error) {
	var cfg Config

	_ = godotenv.Load()

	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
