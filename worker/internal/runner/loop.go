package runner

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/BimaPDev/SignalStack/worker/internal/backoff"
	"github.com/BimaPDev/SignalStack/worker/internal/config"
	"github.com/BimaPDev/SignalStack/worker/internal/processors"
)

type Loop struct {
	db       *sql.DB
	cfg      *config.Config
	registry *processors.Registry
	log      *slog.Logger
}

func New(cfg *config.Config, registry *processors.Registry, log *slog.Logger) (*Loop, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to the database!")
	return &Loop{db: db, cfg: cfg, registry: registry, log: log}, nil
}

func (l *Loop) Run(ctx context.Context) {
	ticker := time.NewTicker(l.cfg.PollInterval)
	defer ticker.Stop()
	defer func() { _ = l.db.Close() }()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.tick(ctx)
		}
	}
}

func (l *Loop) tick(ctx context.Context) {
	job, err := ClaimNext(ctx, l.db, l.cfg.WorkerID)
	if err != nil {
		l.log.Error("failed to claim next job", "error", err)
		return
	}
	if job == nil {
		return
	}

	processor, err := l.registry.Get(job.Type)
	if err != nil {
		l.finalizeJobAsError(ctx, job, fmt.Errorf("no processor registered for type=%q", job.Type))
		return
	}

	startedAt := time.Now().UTC()
	output, procErr := processor.Processor(ctx, job.ID, job.UserID)
	finishedAt := time.Now().UTC()

	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		l.log.Error("begin tx failed", "error", err, "job_id", job.ID)
		return
	}
	defer func() { _ = tx.Rollback() }()

	outputJSON := json.RawMessage("{}")
	if len(output) > 0 {
		outputJSON = output
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO job_results (job_id, output_json, started_at, finished_at)
		VALUES ($1, $2, $3, $4)
	`, job.ID, outputJSON, startedAt, finishedAt)
	if err != nil {
		l.log.Error("insert job_results failed", "error", err, "job_id", job.ID)
		return
	}

	day := startedAt.Format("2006-01-02")

	if procErr == nil {
		l.log.Info("job completed", "job_id", job.ID, "type", job.Type, "user_id", job.UserID)
		_, err = tx.ExecContext(ctx, `
			UPDATE jobs SET status='done', updated_at=now() WHERE id=$1
		`, job.ID)
		if err != nil {
			l.log.Error("update job done failed", "error", err, "job_id", job.ID)
			return
		}
		if err := l.upsertMetric(ctx, tx, job.UserID, day, true); err != nil {
			l.log.Error("metrics update failed", "error", err, "job_id", job.ID)
			return
		}
	} else {
		nextAttempts := job.Attempts + 1
		if nextAttempts < job.MaxAttempts {
			runAt := time.Now().UTC().Add(backoff.NextDelay(nextAttempts))
			_, err = tx.ExecContext(ctx, `
				UPDATE jobs SET status='pending', attempts=$2, run_at=$3, last_error=$4, updated_at=now() WHERE id=$1
			`, job.ID, nextAttempts, runAt, procErr.Error())
		} else {
			_, err = tx.ExecContext(ctx, `
				UPDATE jobs SET status='failed', attempts=$2, last_error=$3, updated_at=now() WHERE id=$1
			`, job.ID, nextAttempts, procErr.Error())
			if err == nil {
				err = l.upsertMetric(ctx, tx, job.UserID, day, false)
			}
		}
		if err != nil {
			l.log.Error("update job failed", "error", err, "job_id", job.ID)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		l.log.Error("commit failed", "error", err, "job_id", job.ID)
	}
}

func (l *Loop) upsertMetric(ctx context.Context, tx *sql.Tx, userID, day string, success bool) error {
	col := "jobs_done"
	if !success {
		col = "jobs_failed"
	}
	_, err := tx.ExecContext(ctx, fmt.Sprintf(`
		WITH existing AS (
			UPDATE metrics_daily SET %s = %s + 1
			WHERE user_id = $1 AND day = $2::date
			RETURNING id
		)
		INSERT INTO metrics_daily (user_id, day, %s)
		SELECT $1, $2::date, 1
		WHERE NOT EXISTS (SELECT 1 FROM existing)
	`, col, col, col), userID, day)
	return err
}

func (l *Loop) finalizeJobAsError(ctx context.Context, job *ClaimedJob, procErr error) {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		l.log.Error("begin tx failed", "error", err, "job_id", job.ID)
		return
	}
	defer func() { _ = tx.Rollback() }()

	now := time.Now().UTC()
	_, err = tx.ExecContext(ctx, `
		INSERT INTO job_results (job_id, output_json, started_at, finished_at)
		VALUES ($1, '{}', $2, $2)
	`, job.ID, now)
	if err != nil {
		l.log.Error("insert job_results failed", "error", err, "job_id", job.ID)
		return
	}

	nextAttempts := job.Attempts + 1
	if nextAttempts < job.MaxAttempts {
		runAt := now.Add(backoff.NextDelay(nextAttempts))
		_, err = tx.ExecContext(ctx, `
			UPDATE jobs SET status='pending', attempts=$2, run_at=$3, last_error=$4, updated_at=now() WHERE id=$1
		`, job.ID, nextAttempts, runAt, procErr.Error())
	} else {
		_, err = tx.ExecContext(ctx, `
			UPDATE jobs SET status='failed', attempts=$2, last_error=$3, updated_at=now() WHERE id=$1
		`, job.ID, nextAttempts, procErr.Error())
	}
	if err != nil {
		l.log.Error("update job failed", "error", err, "job_id", job.ID)
		return
	}

	if err := tx.Commit(); err != nil {
		l.log.Error("commit failed", "error", err, "job_id", job.ID)
	}
}
