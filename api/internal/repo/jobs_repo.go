package repo

// type JobsRepo struct
// - DB *sql.DB

// func NewJobsRepo(db *sql.DB) *JobsRepo
// - return new JobsRepo with db handle

// func (r *JobsRepo) Insert(ctx context.Context, userID string, req model.CreateJobRequest) (*model.CreateJobResponse, error)
// - INSERT INTO jobs (user_id, type, idempotency_key) VALUES (...)
// - handle ON CONFLICT for idempotency_key
// - RETURNING id, status, created_at

// func (r *JobsRepo) List(ctx context.Context, userID string) (*model.ListJobsResponse, error)
// - SELECT * FROM jobs WHERE user_id = $1
// - support optional status filter and pagination

// func (r *JobsRepo) GetByID(ctx context.Context, userID string, id string) (*model.Job, error)
// - SELECT * FROM jobs WHERE id = $1 AND user_id = $2
// - return Job or sql.ErrNoRows
