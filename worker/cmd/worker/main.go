package main

import (
	"log"
	"os"

	"github.com/BimaPDev/SignalStack/api/internal/repo"
	"github.com/joho/godotenv"
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
	godotenv.Load("../.env")
	db, err := repo.Open(os.Getenv("POSTGRES_ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	
}
