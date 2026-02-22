package model

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}

type Event struct {
	ID          string          `json:"id"`
	UserID      string          `json:"user_id"`
	Type        string          `json:"type"`
	PayloadJSON json.RawMessage `json:"payload_json"`
	CreatedAt   time.Time       `json:"created_at"`
}
type Job struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	Type           string     `json:"type"`
	Status         string     `json:"status"`
	RunAt          time.Time  `json:"run_at"`
	Attempts       int        `json:"attempts"`
	MaxAttempts    int        `json:"max_attempts"`
	IdempotencyKey *string    `json:"idempotency_key,omitempty"`
	LockedAt       *time.Time `json:"locked_at,omitempty"`
	LockedBy       *string    `json:"locked_by,omitempty"`
	LastError      *string    `json:"last_error,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
type JobResult struct {
	ID         string          `json:"id"`
	JobID      string          `json:"job_id"`
	OutputJSON json.RawMessage `json:"output_json"`
	StartedAt  time.Time       `json:"started_at"`
	FinishedAt time.Time       `json:"finished_at"`
	CreatedAt  time.Time       `json:"created_at"`
}
type MetricsDaily struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	Day            string `json:"day"`
	EventsReceived int64  `json:"events_received"`
	JobsDone       int64  `json:"jobs_done"`
	JobsFailed     int64  `json:"jobs_failed"`
}
type CreateEventRequest struct {
	Type        string          `json:"type"`
	PayloadJSON json.RawMessage `json:"payload_json"`
}
type CreateEventResponse struct {
	ID        string    `json:"id"`
	CreatedAT time.Time `json:"created_at"`
}
type CreateJobRequest struct {
	Type           string  `json:"type"`
	IdempotencyKey *string `json:"idempotency_key"`
}
type CreateJobResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAT time.Time `json:"created_at"`
}
type ListJobsResponse struct {
	Jobs []Job `json:"jobs"`
}
type AnalyticsSummary struct {
	EventsReceived int64 `json:"events_received"`
	JobsDone       int64 `json:"jobs_done"`
	JobsFailed     int64 `json:"jobs_failed"`
}
type TimeseriesBucket struct {
	Day            string `json:"day"`
	EventsReceived int64  `json:"events_received"`
	JobsDone       int64  `json:"jobs_done"`
	JobsFailed     int64  `json:"jobs_failed"`
}

type TimeseriesResponse struct {
	Buckets []TimeseriesBucket `json:"buckets"`
}
