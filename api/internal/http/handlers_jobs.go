package http

// type JobsHandler struct
// - Repo *repo.JobsRepo
// - Log  *slog.Logger

// func (h *JobsHandler) Create(w http.ResponseWriter, r *http.Request)
// - decode CreateJobRequest from JSON body
// - authenticate request, extract user_id from API key
// - validate type is non-empty
// - call h.Repo.Insert(ctx, userID, req)
// - return 201 with CreateJobResponse as JSON
// - handle idempotency_key conflict (return existing job)

// func (h *JobsHandler) List(w http.ResponseWriter, r *http.Request)
// - authenticate request, extract user_id
// - parse query params: status filter, pagination (limit, offset)
// - call h.Repo.List(ctx, userID)
// - return 200 with ListJobsResponse as JSON

// func (h *JobsHandler) GetByID(w http.ResponseWriter, r *http.Request)
// - authenticate request, extract user_id
// - read {id} path parameter via r.PathValue("id")
// - call h.Repo.GetByID(ctx, userID, id)
// - return 200 with Job as JSON, or 404 if not found
