package repo

import (
	"context"
	"database/sql"

	"github.com/BimaPDev/SignalStack/api/internal/model"
)

// type EventsRepo struct
// - DB *sql.DB

type EventRepo struct {
	DB *sql.DB
}

// func NewEventsRepo(db *sql.DB) *EventsRepo
// - return new EventsRepo with db handle

func NewEventsRepo(db *sql.DB) *EventRepo {
	return &EventRepo{DB: db}
}

// func (r *EventsRepo) Insert(ctx context.Context, userID string, req model.CreateEventRequest) (*model.CreateEventResponse, error)
// - INSERT INTO events (user_id, type, payload_json) VALUES (...)
// - RETURNING id, created_at
// - return CreateEventResponse or error

func (r *EventRepo) Insert(ctx context.Context, userID string, req model.CreateEventRequest) (*model.CreateEventResponse, error) {
	var res model.CreateEventResponse
	err := r.DB.QueryRowContext(ctx, `
		INSERT INTO events (user_id, type, payload_json)
		VALUES ($1,$2,$3)
		RETURNING id, created_at
	`,
		userID, req.Type, req.PayloadJSON,
	).Scan(&res.ID, &res.CreatedAT)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
