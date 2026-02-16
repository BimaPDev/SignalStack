-- 002_indexes.sql

CREATE INDEX IF NOT EXISTS idx_jobs_status_run_at
    ON jobs (status, run_at);

CREATE UNIQUE INDEX IF NOT EXISTS idx_jobs_user_idempotency
    ON jobs (user_id, idempotency_key)
    WHERE idempotency_key IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_events_user_id ON events (user_id);
CREATE INDEX IF NOT EXISTS idx_jobs_user_id ON jobs (user_id);
CREATE INDEX IF NOT EXISTS idx_metrics_daily_user_day ON metrics_daily (user_id, day);
CREATE INDEX IF NOT EXISTS idx_job_results_job_id ON job_results (job_id);
