package http

// type AnalyticsHandler struct
// - Repo *repo.AnalyticsRepo
// - Log  *slog.Logger

// func (h *AnalyticsHandler) Summary(w http.ResponseWriter, r *http.Request)
// - authenticate request, extract user_id
// - parse "from" and "to" query params as dates
// - validate date range
// - call h.Repo.Summary(ctx, userID, from, to)
// - return 200 with AnalyticsSummary as JSON

// func (h *AnalyticsHandler) Timeseries(w http.ResponseWriter, r *http.Request)
// - authenticate request, extract user_id
// - parse "from", "to", "bucket" query params
// - validate date range and bucket value (default "day")
// - call h.Repo.Timeseries(ctx, userID, from, to, bucket)
// - return 200 with TimeseriesResponse as JSON
