package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)


type Config struct {
	DatabaseURL  string
	LogLevel     string
	PollInterval time.Duration
	WorkerID     string
}

func Load() (*Config, error) {
	godotenv.Load(".env", "../.env")
	DatabaseURL := os.Getenv("POSTGRES_ADDR")
	if DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	LogLevel := os.Getenv("LOG_LEVEL")
	if LogLevel == "" {
		LogLevel = "info"
	}
	PollIntervalStr := os.Getenv("POLL_INTERVAL")
	if PollIntervalStr == "" {
		PollIntervalStr = "5s"
	}
	PollInterval, err := time.ParseDuration(PollIntervalStr)
	if err != nil {
		return nil, errors.New("invalid POLL_INTERVAL")
	}
	WorkerID := os.Getenv("WORKER_ID")
	if WorkerID == "" {
		WorkerID = "worker-default"
	}

	return &Config{
		DatabaseURL:  DatabaseURL,
		LogLevel:     LogLevel,
		PollInterval: PollInterval,
		WorkerID:     WorkerID,
	}, nil

}
