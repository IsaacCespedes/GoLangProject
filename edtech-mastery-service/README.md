# Student Progress & Mastery Service

An event-driven Go backend that ingests learning events, computes student mastery and class rollups, and serves teacher dashboards. Built with idempotent ingestion, async workers, and production-style observability.

## Problem

In edtech, teachers need real-time visibility into **mastery** (how well students meet standards) and **progress** (completion, at-risk signals). Raw activity (assignments, submissions, grades) must be aggregated into dashboard-ready metrics without losing events or double-counting. This service demonstrates that pipeline: ingest once, process asynchronously, serve from materialized snapshots.

## Architecture

```
                    POST /events
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  API (chi)                                                   │
│  validate → idempotent insert (events + outbox) → 202         │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────┐    ┌──────────────────────────────────────┐
│  Postgres        │    │  Worker pool (bounded concurrency)    │
│  events (append) │◄───│  claim outbox → update mastery,      │
│  event_outbox    │    │  rollups, risk → mark processed       │
│  student_mastery │    └──────────────────────────────────────┘
│  class_rollups   │
│  risk_flags      │
└──────────────────┘
                         │
                         ▼
              GET /teachers/.../dashboard
              GET /students/.../mastery
              GET /classes/.../students/.../timeline
```

- **Event ingestion**: Append-only `events` table with unique `(source, event_id)` for idempotency. Each new event gets one row in `event_outbox` for the worker.
- **Worker**: Claims pending outbox rows, loads event payload, updates `student_mastery`, `class_rollups`, and `risk_flags` in transactional style. Retries via re-queue on failure; graceful shutdown on SIGINT/SIGTERM.
- **Dashboards**: Read from materialized tables (and events for timeline). Optional Redis caching can be added for dashboard endpoints.

## Event contract

**POST /events** — JSON body:

| Field          | Type     | Required | Description                    |
|----------------|----------|----------|--------------------------------|
| event_id       | string   | yes      | Idempotency key (with source) |
| source         | string   | yes      | System of record               |
| timestamp      | string   | yes      | RFC3339                        |
| student_id     | string   | yes      | Student identifier             |
| class_id       | string   | yes      | Class identifier               |
| assignment_id  | string   | yes*     | *For assignment/submission    |
| standard_ids   | []string | yes*     | *For mastery                  |
| score          | float64  | no       | Present for SUBMISSION_GRADED  |
| rubric_tags    | []string | no       | Optional                       |
| type           | string   | no       | Optional: ASSIGNMENT_ASSIGNED, SUBMISSION_CREATED, SUBMISSION_GRADED (else inferred) |

**Event types** (use `type` or inferred by API):

- `ASSIGNMENT_ASSIGNED` — assignment + standards, no score
- `SUBMISSION_CREATED` — assignment + standards, no score
- `SUBMISSION_GRADED` — assignment + standards + score

## APIs

- **POST /events** — Ingest learning event (idempotent).
- **GET /teachers/{teacherID}/classes/{classID}/dashboard** — Completion rate, average score, at-risk students, recent activity.
- **GET /students/{studentID}/mastery** — Mastery score per standard.
- **GET /classes/{classID}/students/{studentID}/timeline** — Recent event history.
- **GET /metrics** — Prometheus metrics (events ingested, worker latency, dashboard latency, failures).

## At-risk rules

- **missing_submissions**: Student has assignments assigned but no graded submission.
- **score_trend_down**: Last graded score &lt; previous graded score.
- **completion_below_median**: Completion rate below class median.

## Scaling notes (10k+ events/sec)

- **Ingestion**: Partition `events` by `created_at` or hash of `(source, event_id)`; use connection pooling (pgxpool). Idempotency avoids duplicate work on retries.
- **Worker**: Scale horizontally (multiple worker processes claiming from same outbox). Use `FOR UPDATE SKIP LOCKED` (already in place) so workers don’t contend. Consider sharding outbox by `class_id` and dedicated worker pools per shard.
- **Dashboards**: Cache GET responses in Redis with short TTL; use materialized tables to avoid scanning raw events on every request.

## Failure modes

- **Duplicate event**: Same `(source, event_id)` → insert is no-op, no new outbox row; API returns 202 with existing event id.
- **Worker crash**: Outbox row stays `processing` or is reverted to `pending`; another worker can claim and retry. Processing is effectively once per successful commit.
- **DB/queue down**: Ingestion returns 5xx; clients should retry. Worker stops claiming until DB is back.

## Privacy

- Use pseudonymous IDs (`student_id`, `class_id`) in events and storage. No PII in logs; optional correlation IDs for request tracing.

## Run locally

1. **Postgres + Redis** (optional for cache):

   ```bash
   docker compose up -d postgres redis
   ```

2. **Migrations**:

   ```bash
   export DATABASE_URL="postgres://postgres:postgres@localhost:5432/edtech?sslmode=disable"
   migrate -path migrations -database "$DATABASE_URL" up
   ```

3. **API** (port 8080):

   ```bash
   go run ./cmd/api
   ```

4. **Worker** (separate terminal):

   ```bash
   go run ./cmd/worker
   ```

5. **Simulate** (generates events):

   ```bash
   go run ./cmd/simulate
   ```

   Then re-query dashboard; worker will have filled rollups and risk flags.

## Run full stack with Docker

```bash
docker compose up
```

- API: http://localhost:8080  
- Metrics: http://localhost:8080/metrics  
- Postgres: localhost:5432 (user `postgres`, password `postgres`, db `edtech`)  
- Redis: localhost:6379  

## Repo structure

```
/cmd/api        — HTTP API server
/cmd/worker     — Outbox consumer (mastery, rollups, risk)
/cmd/simulate   — Event generator for demo/load test
/internal/domain   — Event and dashboard types
/internal/events   — Validation, ingestion service, processor
/internal/mastery  — Mastery computation from graded events
/internal/risk     — At-risk rules
/internal/rollups  — Class completion/avg score
/internal/dashboard — Dashboard query service
/internal/storage  — Postgres repos (events, outbox, mastery, rollups, risk, timeline)
/internal/queue   — Postgres-backed queue (outbox)
/pkg/logging      — Zerolog setup
/pkg/metrics      — Prometheus counters/histograms
/migrations       — golang-migrate SQL
```

## License

MIT.
