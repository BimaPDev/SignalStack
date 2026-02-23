package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/BimaPDev/SignalStack/api/internal/repo"
)

// type AnalyticsHandler struct
// - Repo *repo.AnalyticsRepo
// - Log  *slog.Logger

type AnalyticsHandler struct {
	Repo *repo.AnalyticsRepo
	Log  *slog.Logger
}

func (h *AnalyticsHandler) Summary(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	res, err := h.Repo.Summary(r.Context(), userID, from, to)
	if err != nil {
		h.Log.Error("summary error:", "err", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// func (h *AnalyticsHandler) Timeseries(w http.ResponseWriter, r *http.Request)
// - authenticate request, extract user_id
// - parse "from", "to", "bucket" query params
// - validate date range and bucket value (default "day")
// - call h.Repo.Timeseries(ctx, userID, from, to, bucket)
// - return 200 with TimeseriesResponse as JSON
