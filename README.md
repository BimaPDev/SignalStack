# SignalStack

A background job queue processing platform built with Go and PostgreSQL. Events come in from HTTP, jobs get queued in the database, and a worker polls and executes them with retry logic and metrics tracking.

## reason to why I made this?

I want to learn :)

## Architecture

```
Client → API (HTTP) → PostgreSQL ← Worker (poller)
```

Two services communicate exclusively through the database — no message broker needed. PostgreSQL acts as the job queue using `FOR UPDATE SKIP LOCKED` for safe concurrent claiming.

- **API** — Receives events, creates jobs, serves status and analytics
- **Worker** — Polls for pending jobs, executes processors, handles retries
- **PostgreSQL** — Shared datastore and job queue

## Getting Started

```bash
# Start PostgreSQL and run migrations
make dev-up

# Run API (port 8080)
make run-api

# Run worker (separate terminal)
make run-worker
```

Or run everything via Docker:

```bash
docker compose up
```

## How Jobs Work

1. POST to `/jobs` → job inserted with `status='pending'`
2. Worker polls every 5s → claims job with `FOR UPDATE SKIP LOCKED`
3. Worker calls the registered processor for that job type
4. On success → `status='done'`, metrics updated
5. On failure → retried with exponential backoff up to `max_attempts`

### Submitting a job manually

```sql
INSERT INTO jobs (user_id, type, status, attempts, max_attempts)
VALUES ('<user_id>', 'example', 'pending', 0, 3);
```

Valid job types must match a registered processor in `worker/cmd/worker/main.go`.

## Processors

Processors are registered in [worker/cmd/worker/main.go](worker/cmd/worker/main.go):

| Type      | Processor          |
| --------- | ------------------ |
| `example` | `ExampleProcessor` |
| `export`  | `ExportProcessor`  |

To add a new job type, implement the `Processor` interface and register it.

## API Endpoints

| Method | Path                    | Description                  |
| ------ | ----------------------- | ---------------------------- |
| POST   | `/events`               | Create an event              |
| POST   | `/jobs`                 | Create a job                 |
| GET    | `/jobs`                 | List jobs (filter by status) |
| GET    | `/jobs/{id}`            | Get job by ID                |
| GET    | `/analytics/summary`    | Total counts for date range  |
| GET    | `/analytics/timeseries` | Day-by-day breakdown         |

All requests require `X-API-Key` header.

## Environment Variables

| Variable        | Required | Default          | Description                  |
| --------------- | -------- | ---------------- | ---------------------------- |
| `POSTGRES_ADDR` | Yes      | —                | PostgreSQL connection string |
| `PORT`          | No       | `8080`           | API server port              |
| `LOG_LEVEL`     | No       | `info`           | Log level                    |
| `POLL_INTERVAL` | No       | `5s`             | Worker poll interval         |
| `WORKER_ID`     | No       | `worker-default` | Worker instance identifier   |

## Project Structure

```
├── api/
│   ├── cmd/api/          # Entry point
│   └── internal/
│       ├── config/       # Config loading
│       ├── http/         # Handlers, router, middleware
│       └── repo/         # Database access layer
├── worker/
│   ├── cmd/worker/       # Entry point + processor registration
│   └── internal/
│       ├── backoff/      # Exponential backoff with jitter
│       ├── config/       # Config loading
│       ├── observability/# Structured logging
│       ├── processors/   # Processor interface + registry
│       └── runner/       # Poll loop + job claiming
└── migrations/           # SQL schema and seed data
```
