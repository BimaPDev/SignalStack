package observability

import (
	"log/slog"
	"os"
)

// func NewLogger() *slog.Logger
// - create a structured JSON logger writing to stdout
// - configure log level from config/environment
// - support hook registration for custom log processing

func NewLogger() *slog.Logger {
	var level slog.Level
	_ = level.UnmarshalText([]byte(os.Getenv("LOG_LEVEL")))

}