# SignalStack — Step-by-Step Guide

This guide walks through every piece of the project, what it does, and why it exists.

---

## Step 1: The Database (migrations/)

Everything starts with your data. Before writing any Go code, you need tables to store things in.

### `001_init.sql` — The Tables

**users** — Every API consumer gets a row here. The `api_key` is how they authenticate.

```sql
CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    api_key     TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

**events** — When a user's app says "something happened", it lands here.

```sql
CREATE TABLE events (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(id),
    type          TEXT NOT NULL,
    payload_json  JSONB NOT NULL DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

- `type` — what kind of event: `"payment.completed"`, `"user.signup"`, etc.
- `payload_json` — any extra data as JSON. `JSONB` means Postgres can index and query inside it.

**jobs** — Work that needs to be done. A job is created in response to an event.

```sql
CREATE TABLE jobs (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID NOT NULL REFERENCES users(id),
    type             TEXT NOT NULL,
    status           TEXT NOT NULL DEFAULT 'pending',
    run_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    attempts         INT NOT NULL DEFAULT 0,
    max_attempts     INT NOT NULL DEFAULT 3,
    idempotency_key  TEXT,
    locked_at        TIMESTAMPTZ,
    locked_by        TEXT,
    last_error       TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

Key columns to understand:
- `status` — lifecycle: `pending` → `running` → `done` or `failed`
- `run_at` — when should this job run? Allows delayed/scheduled jobs
- `attempts` / `max_attempts` — retry tracking. Default 3 tries before giving up
- `idempotency_key` — prevents duplicate jobs (e.g., double-clicking "submit")
- `locked_by` — which worker claimed this job (prevents two workers running the same job)

**job_results** — After a job runs, the output is stored here.

```sql
CREATE TABLE job_results (
    job_id       UUID NOT NULL REFERENCES jobs(id),
    output_json  JSONB NOT NULL DEFAULT '{}',
    started_at   TIMESTAMPTZ NOT NULL,
    finished_at  TIMESTAMPTZ NOT NULL
);
```

**metrics_daily** — Aggregated counts per user per day. Used by the analytics endpoints.

```sql
CREATE TABLE metrics_daily (
    user_id         UUID NOT NULL REFERENCES users(id),
    day             DATE NOT NULL,
    events_received BIGINT NOT NULL DEFAULT 0,
    jobs_done       BIGINT NOT NULL DEFAULT 0,
    jobs_failed     BIGINT NOT NULL DEFAULT 0
);
```

### `002_indexes.sql` — Making Queries Fast

```sql
-- The worker polls for pending jobs ordered by run_at.
-- Without this index, it would scan every row in the jobs table.
CREATE INDEX idx_jobs_status_run_at ON jobs (status, run_at);

-- Ensures one idempotency_key per user. The WHERE clause
-- makes it a partial index — only rows with a key are indexed.
CREATE UNIQUE INDEX idx_jobs_user_idempotency
    ON jobs (user_id, idempotency_key)
    WHERE idempotency_key IS NOT NULL;
```

**Rule of thumb:** if you query by a column in a WHERE clause, it probably needs an index.

---

## Step 2: Configuration (config/)

Before your app can do anything, it needs to know: where's the database? what port to listen on?

```go
type Config struct {
    DatabaseURL string   // connection string for Postgres
    Port        string   // HTTP port (default "8080")
    LogLevel    string   // "debug", "info", "warn", "error"
}
```

`Load()` reads from **environment variables**. This is the 12-factor app approach — config lives outside the code so the same binary works in dev, staging, and production.

---

## Step 3: Database Connection (repo/db.go)

```go
func Open(databaseURL string) (*sql.DB, error) {
    db, err := sql.Open("postgres", databaseURL)
    // ...
    err = db.Ping()  // actually try to connect
    // ...
    return db, nil
}
```

**Why `sql.Open` + `Ping`?** `sql.Open` only validates the connection string — it doesn't actually connect. `Ping` forces a real connection attempt so you fail fast at startup instead of on the first request.

---

## Step 4: Models (model/types.go)

Two categories of types:

**Entities** — mirror your database tables exactly:

```go
type Job struct {
    ID             string
    UserID         string
    Type           string
    Status         string     // "pending", "running", "done", "failed"
    Attempts       int
    MaxAttempts    int
    // ...
}
```

**DTOs (Data Transfer Objects)** — what the API sends/receives:

```go
type CreateEventRequest struct {
    Type        string
    PayloadJSON json.RawMessage
}

type CreateEventResponse struct {
    ID        string
    CreatedAt time.Time
}
```

**Why separate them?** Your database row has 12 fields. Your API response might only need 3. DTOs let you control exactly what goes in and out without exposing internal fields like `locked_by` or `last_error`.

---

## Step 5: Repositories (repo/)

Repos are the **only** layer that talks to the database. Each repo owns one table.

### events_repo.go

```go
type EventsRepo struct {
    DB *sql.DB
}

func (r *EventsRepo) Insert(ctx, userID, req) (*CreateEventResponse, error) {
    // INSERT INTO events (user_id, type, payload_json) VALUES (...)
    // RETURNING id, created_at
}
```

### jobs_repo.go

```go
func (r *JobsRepo) Insert(...)    // create a job, handle idempotency conflicts
func (r *JobsRepo) List(...)      // fetch jobs with optional filters
func (r *JobsRepo) GetByID(...)   // fetch one job, return error if not found
```

### analytics_repo.go

```go
func (r *AnalyticsRepo) Summary(...)     // SUM of metrics over a date range
func (r *AnalyticsRepo) Timeseries(...)  // metrics broken down by day
```

**Why this pattern?** Handlers never write SQL. If you need to switch databases, change the query format, or add caching — you only touch the repo layer. Everything above it stays the same.

---

## Step 6: HTTP Handlers (http/handlers_*.go)

Handlers are the **glue** between HTTP and your business logic. Each one follows the same pattern:

```
1. Parse the request (decode JSON, read URL params)
2. Validate the input (is the type field empty?)
3. Call the repo (do the actual work)
4. Return a response (JSON with the right status code)
```

### EventHandler

```go
type EventHandler struct {
    repo *repo.EventRepo    // database access
    Log  *slog.Logger       // structured logging
}
```

The handler **doesn't know** how the database works. It just calls `h.repo.Insert(...)` and gets back a response. This is dependency injection — the handler receives what it needs through its struct fields.

### JobsHandler

Same pattern but with three methods:
- `Create` — make a new job
- `List` — fetch jobs with filters
- `GetByID` — fetch one job

### AnalyticsHandler

- `Summary` — total counts for a date range
- `Timeseries` — counts broken down by day

---

## Step 7: Middleware (http/middleware.go)

Middleware wraps every request. It runs **before** and **after** your handler.

### RequestIDMiddleware

```
Request comes in
  → Check for X-Request-ID header
  → If missing, generate a UUID
  → Store it in the request context
  → Pass to next handler
  → Set X-Request-ID on the response
```

**Why?** When debugging in production, a request ID lets you trace a single request across all your logs.

### LoggingMiddleware

```
Request comes in
  → Record start time
  → Pass to next handler
  → Calculate duration
  → Log: method, path, status, duration_ms, request_id
```

Every request gets a structured log line. When something breaks, you can search your logs by path, status code, or request ID.

---

## Step 8: Router (http/router.go)

The router is where everything gets **wired together**. It maps URL paths to handler methods and wraps them with middleware.

```
GET  /health                  → handleHealth (simple status check)
POST /events                  → EventHandler.Create
POST /jobs                    → JobsHandler.Create
GET  /jobs                    → JobsHandler.List
GET  /jobs/{id}               → JobsHandler.GetByID
GET  /analytics/summary       → AnalyticsHandler.Summary
GET  /analytics/timeseries    → AnalyticsHandler.Timeseries
```

The router function needs to **receive dependencies** (database, logger) so it can create the handler structs.

---

## Step 9: main.go — Where It All Comes Together

```
1. Load config from environment
2. Open database connection
3. Create logger
4. Create router (pass in db + logger)
5. Start HTTP server
6. Listen for shutdown signals
```

This is the **composition root** — the one place that creates all dependencies and wires them together. No other file creates database connections or loggers. Everything flows from here.

---

## Step 10: The Worker (worker/)

The worker is a separate binary that runs alongside the API. It doesn't serve HTTP — it polls the database for pending jobs.

### The Loop (runner/loop.go)

```
Every N seconds:
  1. Claim the next pending job (atomically)
  2. Look up the right processor for the job type
  3. Run the processor
  4. Save the result
  5. Update the job status (done or failed)
  6. Update daily metrics
```

### Claiming Jobs (runner/claim_stub.go)

This is the trickiest part. Multiple workers might run at once. You need to ensure **only one worker** picks up each job.

```sql
UPDATE jobs SET status='running', locked_by=$1
WHERE id = (
    SELECT id FROM jobs
    WHERE status='pending' AND run_at <= now()
    ORDER BY run_at
    FOR UPDATE SKIP LOCKED
    LIMIT 1
)
RETURNING id, type, user_id
```

**`FOR UPDATE SKIP LOCKED`** — this is PostgreSQL magic. It locks the row so no other worker can claim it, and `SKIP LOCKED` means if a row is already locked, skip it instead of waiting. This gives you a **safe, concurrent job queue** with just SQL.

### Processors (processors/)

```go
type Processor interface {
    Process(ctx context.Context, jobID, userID string) ([]byte, error)
}
```

Each job type gets its own processor. The registry maps type names to implementations:

```go
registry.Register("send_email", &EmailProcessor{})
registry.Register("generate_report", &ReportProcessor{})
```

When the worker claims a job with `type="send_email"`, it looks up the processor and calls `Process()`.

### Backoff (backoff/backoff.go)

When a job fails, you don't retry immediately. You wait longer each time:

```
Attempt 1 fails → wait ~2 seconds
Attempt 2 fails → wait ~4 seconds
Attempt 3 fails → give up, mark as "failed"
```

**Jitter** adds randomness so that if 100 jobs all fail at once, they don't all retry at the exact same moment (which would overload whatever they're calling).

---

## How The Pieces Connect

```
                    ┌─────────────┐
                    │   Client    │
                    └──────┬──────┘
                           │ HTTP
                    ┌──────▼──────┐
                    │   API       │
                    │  (chi)      │
                    └──────┬──────┘
                           │ SQL
                    ┌──────▼──────┐
                    │  PostgreSQL │
                    └──────▲──────┘
                           │ SQL (poll)
                    ┌──────┴──────┐
                    │   Worker    │
                    │  (loop)     │
                    └─────────────┘
```

The API and Worker **never talk to each other directly**. They communicate through the database. The API writes jobs, the worker reads them. This is simple, reliable, and easy to scale — just run more workers.

---

## What To Build Next

Now that you understand every piece, here's the order to implement:

1. **Wire up the router** — connect your real handlers to routes
2. **Implement the handler methods** — replace the pseudocode comments with real Go
3. **Implement the repo methods** — write the actual SQL queries
4. **Add auth middleware** — validate API keys, inject user_id into context
5. **Write tests** — start with repos, then handlers
6. **Implement a real processor** — replace the example with something useful
