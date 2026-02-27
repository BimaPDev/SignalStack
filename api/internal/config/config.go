package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// type Config struct
// - DatabaseURL string
// - Port        string
// - LogLevel    string

type Config struct {
	DatabaseURL string
	Port string
	LogLevel string
}
// func Load() (*Config, error)
// - read DATABASE_URL, PORT, LOG_LEVEL from environment
// - return error if DATABASE_URL is missing
// - default PORT to "8080", LOG_LEVEL to "info"

func Load() (*Config, error) {
	godotenv.Load("../.env")
	DatabaseURL := os.Getenv("POSTGRES_ADDR")
	if DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
	    port = "8080"
	}
	LogLevel := os.Getenv("LOG_LEVEL")
	if LogLevel == "" {
		LogLevel ="info"
	}

	return &Config{
		DatabaseURL: DatabaseURL,
		Port:	port,
		LogLevel: LogLevel,
	}, nil
}
