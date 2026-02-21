# SignalStack — Architecture Overview

## What Is SignalStack?

A **webhook/event processing platform** — a simplified version of services like Zapier, AWS EventBridge, or a job queue system.

## The Flow

```
Client sends event → API receives it → Job gets created → Worker picks it up → Worker executes it
```

### Real-World Example

1. A user's app sends: "a payment happened" (`POST /events`)
2. SignalStack creates a job: "send a confirmation email"
3. The worker picks up that job, runs it, retries if it fails
4. Analytics track how many events/jobs succeeded/failed per day

## Who Would Use This?

Any developer who needs **reliable background processing**:

- Sending emails/notifications after an event
- Processing payments asynchronously
- Running scheduled tasks with retry logic

## Components

| Component      | Purpose                                                 |
| -------------- | ------------------------------------------------------- |
| **API**        | Receives events, creates jobs, serves analytics         |
| **Worker**     | Polls for pending jobs, executes them, handles retries  |
| **PostgreSQL** | Stores everything, used as the job queue (`SKIP LOCKED`)|
| **Metrics**    | Tracks daily counts of events and job outcomes          |

## Project Structure

```
SignalStack/
├── api/                    # HTTP API server
│   ├── cmd/api/            # Entry point (main.go)
│   └── internal/
│       ├── config/         # Environment config loader
│       ├── http/           # Router, handlers, middleware
│       ├── model/          # Structs and DTOs
│       ├── observability/  # Structured logging
│       └── repo/           # Database access layer
├── worker/                 # Background job processor
│   ├── cmd/worker/         # Entry point (main.go)
│   └── internal/
│       ├── backoff/        # Exponential backoff with jitter
│       ├── config/         # Environment config loader
│       ├── observability/  # Structured logging
│       ├── processors/     # Job type implementations
│       └── runner/         # Polling loop and job claiming
├── migrations/             # SQL schema and indexes
├── docker-compose.yml      # Local dev environment
└── Makefile                # Build and run commands
```

## Request Lifecycle

```
HTTP Request
  → Router (chi)
    → Middleware (request ID, logging)
      → Handler (validates input, calls repo)
        → Repo (executes SQL, returns data)
          → Handler (formats JSON response)
            → HTTP Response
```

## API Endpoints

| Method | Path                      | Description                  |
| ------ | ------------------------- | ---------------------------- |
| GET    | `/health`                 | Health check                 |
| POST   | `/events`                 | Create a new event           |
| POST   | `/jobs`                   | Create a new job             |
| GET    | `/jobs`                   | List jobs (with filters)     |
| GET    | `/jobs/{id}`              | Get a single job by ID       |
| GET    | `/analytics/summary`      | Aggregated metrics           |
| GET    | `/analytics/timeseries`   | Time-bucketed metrics        |

## Key Concepts

- **Dependency Injection**: Handlers receive their dependencies (repos, logger) via structs — no globals
- **Idempotency**: Jobs use an `idempotency_key` to prevent duplicate processing
- **SKIP LOCKED**: Worker uses PostgreSQL row-level locking for safe concurrent job claiming
- **Exponential Backoff**: Failed jobs are retried with increasing delays + jitter
- **Graceful Shutdown**: Worker listens for SIGINT/SIGTERM and finishes in-flight work
