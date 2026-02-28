package main

import (
	"context"
	"fmt"

	"github.com/BimaPDev/SignalStack/worker/internal/config"
	"github.com/BimaPDev/SignalStack/worker/internal/observability"
	"github.com/BimaPDev/SignalStack/worker/internal/processors"
	"github.com/BimaPDev/SignalStack/worker/internal/runner"
	_ "github.com/lib/pq"
)

// func main
// - load config from environment via config.Load()
// - create structured logger via observability.NewLogger()
// - create processor registry via processors.NewRegistry()
// - register processor implementations on registry
// - create runner loop via runner.New(cfg, registry, logger)
// - set up context with cancel for graceful shutdown
// - listen for SIGINT/SIGTERM, cancel context on signal
// - call loop.Run(ctx) — blocks until context cancelled

func main() {
	logger := observability.NewLogger()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		return
	}

	registry := processors.NewRegistry()
	registry.Register("example", &processors.ExampleProcessor{})
	registry.Register("export", &processors.ExportProcessor{})

	loop, err := runner.New(cfg, registry, logger)
	if err != nil {
		logger.Error("Failed to start loop", "error", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Print("Workers Online and ready to go \n")
	loop.Run(ctx)
}
