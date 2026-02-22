package http

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/BimaPDev/SignalStack/api/internal/repo"
	"github.com/go-chi/chi/v5"
)

// func NewRouter(db *sql.DB, log *slog.Logger) http.Handler
// - create http.NewServeMux()
// - instantiate EventsHandler, JobsHandler, AnalyticsHandler with repos
// - register routes:
//     GET  /health                    -> handleHealth
//     POST /events                    -> EventsHandler.Create
//     POST /jobs                      -> JobsHandler.Create
//     GET  /jobs                      -> JobsHandler.List
//     GET  /jobs/{id}                 -> JobsHandler.GetByID
//     GET  /analytics/summary         -> AnalyticsHandler.Summary
//     GET  /analytics/timeseries      -> AnalyticsHandler.Timeseries
// - wrap mux with RequestIDMiddleware and LoggingMiddleware
// - return wrapped handler

// func handleHealth(w http.ResponseWriter, r *http.Request)
// - respond 200 with {"status":"ok"}

func NewRouter(db *sql.DB, log *slog.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(RequestIDMiddleware)
	r.Use(LoggingMiddleware(log))
	r.Get("/health", handleHealth)
	EV := &EventHandler{repo: &repo.EventRepo{DB: db},
		Log: log}
	r.Post("/events", EV.Create)
	return r
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"status":"ok"}`))
}
