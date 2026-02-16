package repo

// type AnalyticsRepo struct
// - DB *sql.DB

// func NewAnalyticsRepo(db *sql.DB) *AnalyticsRepo
// - return new AnalyticsRepo with db handle

// func (r *AnalyticsRepo) Summary(ctx context.Context, userID, from, to string) (*model.AnalyticsSummary, error)
// - SELECT SUM(events_received), SUM(jobs_done), SUM(jobs_failed)
//   FROM metrics_daily WHERE user_id = $1 AND day BETWEEN $2 AND $3
// - return AnalyticsSummary or error

// func (r *AnalyticsRepo) Timeseries(ctx context.Context, userID, from, to, bucket string) (*model.TimeseriesResponse, error)
// - SELECT day, events_received, jobs_done, jobs_failed
//   FROM metrics_daily WHERE user_id = $1 AND day BETWEEN $2 AND $3
//   ORDER BY day
// - return TimeseriesResponse or error
