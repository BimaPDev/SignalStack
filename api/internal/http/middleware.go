package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"database/sql"

	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "request_id"
const userIDKey contextKey = "user_id"

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoggingMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recordTime := time.Now()
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)
			log.Info("request", "method", r.Method, "path", r.URL.Path, "status", rw.status, "requestID", r.Context().Value(RequestIDKey), "timestamp", time.Since(recordTime).Milliseconds())
		})
	}
}

func AuthMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	// 1. read X-API-Key header
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				http.Error(w, `{"error":"APIKey Not Found"}`, http.StatusUnauthorized)
				return
			}
			// 2. query DB for user
			var userID string
			err := db.QueryRowContext(r.Context(), `
    		SELECT id FROM users WHERE api_key = $1
		`, apiKey).Scan(&userID)
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			if err != nil {
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
				return
			}
			// 3. 401 if not found
			// 4. inject user_id into context
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
