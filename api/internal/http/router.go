package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

func HandlerRouter()  http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", handleHealth)
	return r
}

func handleHealth(w http.ResponseWriter, r *http.Request){
	w.Write([]byte(`{"status":"ok"}`))
}

