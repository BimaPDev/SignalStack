package model

// --- Entities (mirror DB tables) ---

// type User struct
// - ID        string    (uuid)
// - APIKey    string
// - CreatedAt time.Time

// type Event struct
// - ID          string          (uuid)
// - UserID      string          (fk -> users)
// - Type        string
// - PayloadJSON json.RawMessage
// - CreatedAt   time.Time

// type Job struct
// - ID             string     (uuid)
// - UserID         string     (fk -> users)
// - Type           string
// - Status         string     (pending | running | done | failed)
// - RunAt          time.Time
// - Attempts       int
// - MaxAttempts    int
// - IdempotencyKey *string
// - LockedAt       *time.Time
// - LockedBy       *string
// - LastError      *string
// - CreatedAt      time.Time
// - UpdatedAt      time.Time

// type JobResult struct
// - ID         string          (uuid)
// - JobID      string          (fk -> jobs)
// - OutputJSON json.RawMessage
// - StartedAt  time.Time
// - FinishedAt time.Time
// - CreatedAt  time.Time

// type MetricsDaily struct
// - ID             string (uuid)
// - UserID         string (fk -> users)
// - Day            string (YYYY-MM-DD)
// - EventsReceived int64
// - JobsDone       int64
// - JobsFailed     int64

// --- Request/Response DTOs ---

// type CreateEventRequest struct
// - Type        string
// - PayloadJSON json.RawMessage

// type CreateEventResponse struct
// - ID        string
// - CreatedAt time.Time

// type CreateJobRequest struct
// - Type           string
// - IdempotencyKey *string

// type CreateJobResponse struct
// - ID        string
// - Status    string
// - CreatedAt time.Time

// type ListJobsResponse struct
// - Jobs []Job

// type AnalyticsSummary struct
// - EventsReceived int64
// - JobsDone       int64
// - JobsFailed     int64

// type TimeseriesBucket struct
// - Day            string
// - EventsReceived int64
// - JobsDone       int64
// - JobsFailed     int64

// type TimeseriesResponse struct
// - Buckets []TimeseriesBucket
