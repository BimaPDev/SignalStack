package http

// type contextKey string
// const RequestIDKey contextKey = "request_id"

// func RequestIDMiddleware(next http.Handler) http.Handler
// - read X-Request-ID header from request; if empty, generate uuid
// - store request ID in context via context.WithValue
// - set X-Request-ID response header
// - call next.ServeHTTP

// func LoggingMiddleware(next http.Handler, log *slog.Logger) http.Handler
// - record request start time
// - wrap ResponseWriter to capture status code
// - call next.ServeHTTP
// - log method, path, status, duration_ms, request_id as structured fields
