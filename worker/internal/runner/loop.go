package runner

// type Loop struct
// - db       *sql.DB
// - cfg      *config.Config
// - registry *processors.Registry
// - log      *slog.Logger

// func New(cfg *config.Config, registry *processors.Registry, log *slog.Logger) (*Loop, error)
// - open postgres connection using cfg.DatabaseURL
// - ping to verify connectivity
// - return *Loop or error

// func (l *Loop) Run(ctx context.Context)
// - start time.Ticker with cfg.PollInterval
// - loop: select on ctx.Done() (return) or ticker.C (call tick)
// - defer ticker.Stop() and db.Close()

// func (l *Loop) tick(ctx context.Context)
// - call ClaimNext(ctx, l.db, l.cfg.WorkerID) to atomically claim next runnable job
// - if no job available, return (nothing to do this tick)
// - look up processor via l.registry.Get(job.Type)
// - call processor.Process(ctx, job.ID, job.UserID)
// - INSERT into job_results with output, started_at, finished_at
// - on success: UPDATE jobs SET status='done'
// - on failure with attempts < max_attempts: UPDATE jobs SET status='pending', run_at=now()+backoff
// - on failure with attempts >= max_attempts: UPDATE jobs SET status='failed', last_error=err
// - UPDATE metrics_daily: increment jobs_done or jobs_failed for user+day
