package runner

// type ClaimedJob struct
// - ID     string
// - Type   string
// - UserID string

// func ClaimNext(ctx context.Context, db *sql.DB, workerID string) (*ClaimedJob, error)
// - execute atomic claim SQL:
//     UPDATE jobs SET status='running', locked_by=$1, locked_at=now(), attempts=attempts+1
//     WHERE id = (
//       SELECT id FROM jobs
//       WHERE status='pending' AND run_at <= now()
//       ORDER BY run_at
//       FOR UPDATE SKIP LOCKED
//       LIMIT 1
//     )
//     RETURNING id, type, user_id
// - if no rows returned, return (nil, nil) meaning no work available
// - otherwise return *ClaimedJob
