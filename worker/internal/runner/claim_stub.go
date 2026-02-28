package runner

import (
	"context"
	"database/sql"
)

type ClaimedJob struct {
	ID          string
	Type        string
	UserID      string
	Attempts    int
	MaxAttempts int
}

func ClaimNext(ctx context.Context, db *sql.DB, workerID string) (*ClaimedJob, error) {
	row := db.QueryRowContext(ctx, `
		UPDATE jobs SET status='running', locked_by=$1, locked_at=now(), attempts=attempts+1
		WHERE id = (
			SELECT id FROM jobs
			WHERE status='pending' AND run_at <= now()
			ORDER BY run_at
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		RETURNING id, type, user_id, attempts, max_attempts
	`, workerID)

	var job ClaimedJob
	err := row.Scan(&job.ID, &job.Type, &job.UserID, &job.Attempts, &job.MaxAttempts)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &job, nil
}
