package config

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
