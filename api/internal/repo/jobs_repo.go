package repo

import (
	"context"
	"database/sql"

	"github.com/BimaPDev/SignalStack/api/internal/model"
)

type JobsRepo struct {
	DB *sql.DB
}

func NewJobsRepo(db *sql.DB) *JobsRepo {
	return &JobsRepo{DB: db}
}

// func (r *JobsRepo) Insert(ctx context.Context, userID string, req model.CreateJobRequest) (*model.CreateJobResponse, error)
// - INSERT INTO jobs (user_id, type, idempotency_key) VALUES (...)
// - handle ON CONFLICT for idempotency_key
// - RETURNING id, status, created_at

func (r *JobsRepo) Insert(ctx context.Context, userID string, req model.CreateJobRequest) (*model.CreateJobResponse, error) {
	var res model.CreateJobResponse
	err := r.DB.QueryRowContext(ctx, `
		INSERT INTO jobs (user_id, type, idempotency_key)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, idempotency_key) WHERE idempotency_key IS NOT NULL
		DO NOTHING
		RETURNING id, status, created_at
		`,
		userID, req.Type, req.IdempotencyKey,
	).Scan(&res.ID, &res.Status, &res.CreatedAT)
	if err == sql.ErrNoRows {
		err := r.DB.QueryRowContext(ctx, `
			
		`
	} else if err != nil {
		return nil, err
	}
	return &res, nil
}

// func (r *JobsRepo) List(ctx context.Context, userID string) (*model.ListJobsResponse, error)
// - SELECT * FROM jobs WHERE user_id = $1
// - support optional status filter and pagination

// func (r *JobsRepo) GetByID(ctx context.Context, userID string, id string) (*model.Job, error)
// - SELECT * FROM jobs WHERE id = $1 AND user_id = $2
// - return Job or sql.ErrNoRows
