package http

import (
	"log/slog"

	"github.com/BimaPDev/SignalStack/api/internal/repo"
)

// type EventsHandler struct
// - Repo *repo.EventsRepo
// - Log  *slog.Logger

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
