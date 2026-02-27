package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// type Config struct
// - DatabaseURL  string
// - LogLevel     string
// - PollInterval time.Duration
// - WorkerID     string

// func Load() (*Config, error)
// - read DATABASE_URL, LOG_LEVEL, POLL_INTERVAL, WORKER_ID from environment
// - return error if DATABASE_URL is missing
// - default POLL_INTERVAL to 5s, WORKER_ID to "worker-default", LOG_LEVEL to "info"
// - parse POLL_INTERVAL as time.Duration

type Config struct {
	DatabaseURL  string
	LogLevel     string
	PollInterval time.Duration
	WorkerID     string
}

func Load() (*Config, error) {
	godotenv.Load("../.env")
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
