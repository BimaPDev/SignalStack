package repo

// type EventsRepo struct
// - DB *sql.DB

// func NewEventsRepo(db *sql.DB) *EventsRepo
// - return new EventsRepo with db handle

// func (r *EventsRepo) Insert(ctx context.Context, userID string, req model.CreateEventRequest) (*model.CreateEventResponse, error)
// - INSERT INTO events (user_id, type, payload_json) VALUES (...)
// - RETURNING id, created_at
// - return CreateEventResponse or error
