package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/BimaPDev/SignalStack/api/internal/model"
	"github.com/BimaPDev/SignalStack/api/internal/repo"
)

type EventHandler struct {
	repo *repo.EventRepo
	Log  *slog.Logger
}

// func (h *EventsHandler) Create(w http.ResponseWriter, r *http.Request)
// - decode CreateEventRequest from JSON body
// - authenticate request, extract user_id from API key
// - validate type is non-empty
// - call h.Repo.Insert(ctx, userID, req)
// - return 201 with CreateEventResponse as JSON
// - on error: return appropriate HTTP status + error JSON

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.Type == "" {
		http.Error(w, `{"error":"req.Type not found"}`, http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	res, err := h.repo.Insert(r.Context(), userID, req)
	if err != nil {
		h.Log.Error("insert failed", "err", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
