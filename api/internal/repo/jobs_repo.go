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
			SELECT id, status, created_at
			FROM jobs
			WHERE user_id = $1 AND idempotency_key = $2
		`,
			userID, req.IdempotencyKey,
		).Scan(&res.ID, &res.Status, &res.CreatedAT)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *JobsRepo) List(ctx context.Context, userID string, status string, limit int, offset int) (*model.ListJobsResponse, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, user_id, type, status, run_at, attempts, max_attempts,
       	idempotency_key, locked_at, locked_by, last_error, created_at, updated_at
		FROM jobs
		WHERE user_id = $1
		AND ($2 = '' OR status = $2)
		LIMIT $3 OFFSET $4
	`, userID, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []model.Job
	for rows.Next() {
		var job model.Job
		if err := rows.Scan(&job.ID, &job.UserID, &job.Type, &job.Status, &job.RunAt, &job.Attempts, &job.MaxAttempts, &job.IdempotencyKey, &job.LockedAt, &job.LockedBy, &job.LastError, &job.CreatedAt, &job.UpdatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return &model.ListJobsResponse{Jobs: jobs}, nil

}
func (r *JobsRepo) GetByID(ctx context.Context, userID string, id string) (*model.Job, error) {
	var job model.Job
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, user_id, type, status, run_at, attempts, max_attempts,
       	idempotency_key, locked_at, locked_by, last_error, created_at, updated_at
		FROM jobs
		WHERE id = $1 AND user_id = $2`, id, userID).Scan(&job.ID, &job.UserID, &job.Type, &job.Status, &job.RunAt, &job.Attempts, &job.MaxAttempts, &job.IdempotencyKey, &job.LockedAt, &job.LockedBy, &job.LastError, &job.CreatedAt, &job.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &job, nil
}
