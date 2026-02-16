package main

// func main
// - load config from environment via config.Load()
// - create structured logger via observability.NewLogger()
// - create processor registry via processors.NewRegistry()
// - register processor implementations on registry
// - create runner loop via runner.New(cfg, registry, logger)
// - set up context with cancel for graceful shutdown
// - listen for SIGINT/SIGTERM, cancel context on signal
// - call loop.Run(ctx) — blocks until context cancelled
