package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/BimaPDev/SignalStack/api/internal/config"
	internalhttp "github.com/BimaPDev/SignalStack/api/internal/http"
	"github.com/BimaPDev/SignalStack/api/internal/repo"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	db, err := repo.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	router := internalhttp.NewRouter(db, slog.Default())
	addr := fmt.Sprintf(":%s", cfg.Port)
	slog.Default().Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
