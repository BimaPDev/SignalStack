package repo

import (
	"context"
	"database/sql"

	"github.com/BimaPDev/SignalStack/api/internal/model"
)

type AnalyticsRepo struct {
	DB *sql.DB
}


func NewAnalyticsRepo(db *sql.DB) *AnalyticsRepo {
	return &AnalyticsRepo{DB: db}
}


func (r *AnalyticsRepo) Summary(ctx context.Context, userID string, from string, to string) (*model.AnalyticsSummary, error) {
	var sum model.AnalyticsSummary
	err := r.DB.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(events_received),0), COALESCE(SUM(jobs_done),0), COALESCE(SUM(jobs_failed),0)
		from metrics_daily WHERE user_id = $1 AND day BETWEEN $2 AND $3
	`,
		userID, from, to,
	).Scan(&sum.EventsReceived, &sum.JobsDone, &sum.JobsFailed)
	if err != nil {
		return nil, err
	}
	return &sum, nil
}

// func (r *AnalyticsRepo) Timeseries(ctx context.Context, userID, from, to, bucket string) (*model.TimeseriesResponse, error)
// - SELECT day, events_received, jobs_done, jobs_failed
//   FROM metrics_daily WHERE user_id = $1 AND day BETWEEN $2 AND $3
//   ORDER BY day
// - return TimeseriesResponse or error
