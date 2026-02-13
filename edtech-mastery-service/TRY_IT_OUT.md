# How to Try Out the Student Progress & Mastery Service

This guide walks you through running the app and exercising its APIs so you can see events, dashboards, and metrics in action.

---

## Prerequisites

- **Docker** and **Docker Compose** (for Postgres, Redis, and optionally the full stack)
- **Go 1.22+** (if you run the API and worker locally)
- **golang-migrate** CLI (only if you run migrations yourself; not needed when using `docker compose up`)

---

## Option A: Run Everything with Docker (Easiest)

From the `edtech-mastery-service` directory:

```bash
docker compose up
```

This starts Postgres, runs migrations, then the API and worker. When you see the API and worker logs, the stack is ready.

- **API**: http://localhost:8080  
- **Metrics**: http://localhost:8080/metrics  
- **Postgres**: localhost:5432 (user `postgres`, password `postgres`, db `edtech`)  
- **Redis**: localhost:6379  

Then skip to [Generate demo data](#generate-demo-data).

---

## Option B: Run Locally (API + Worker + Simulate)

Useful if you want to change code and restart quickly.

### 1. Start Postgres (and optionally Redis)

```bash
docker compose up -d postgres redis
```

### 2. Run migrations

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/edtech?sslmode=disable"
migrate -path migrations -database "$DATABASE_URL" up
```

*(Install migrate: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0`)*

### 3. Start the API (terminal 1)

```bash
go run ./cmd/api
```

Leave it running. Default port: **8080**.

### 4. Start the worker (terminal 2)

```bash
go run ./cmd/worker
```

Leave it running. It processes events from the outbox and updates mastery, rollups, and risk flags.

### 5. Generate demo data (terminal 3)

```bash
go run ./cmd/simulate
```

This sends many events to the API (assignments assigned, submissions created, submissions graded) for **teacher-1**, **class-1**, and 20 students by default. You can override:

```bash
STUDENTS=5 ASSIGNMENTS=3 go run ./cmd/simulate
```

After the simulator finishes, the worker will process the outbox. Give it a few seconds, then use the API calls below to see dashboards and mastery.

---

## Try the APIs

Base URL: **http://localhost:8080** (use `API_URL` if you changed it).

### 1. Ingest a single event (POST /events)

```bash
curl -s -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "ev-manual-1",
    "source": "manual",
    "timestamp": "2025-02-13T12:00:00Z",
    "student_id": "student-1",
    "class_id": "class-1",
    "assignment_id": "assign-1",
    "standard_ids": ["std-math-1", "std-math-2"],
    "type": "ASSIGNMENT_ASSIGNED"
  }'
```

Expect: **202 Accepted**. Sending the same `event_id` + `source` again is idempotent (no duplicate outbox row).

### 2. Teacher dashboard (GET)

After running the simulator (and letting the worker run), the dashboard for the demo class is:

```bash
curl -s http://localhost:8080/teachers/teacher-1/classes/class-1/dashboard | jq
```

You should see completion rate, average score, at-risk students, and recent activity. *(Install `jq` for pretty JSON, or omit it.)*

### 3. Student mastery (GET)

```bash
curl -s http://localhost:8080/students/student-1/mastery | jq
```

Shows mastery per standard for that student.

### 4. Student timeline (GET)

```bash
curl -s "http://localhost:8080/classes/class-1/students/student-1/timeline" | jq
```

Recent event history for that student in that class.

### 5. Prometheus metrics (GET)

```bash
curl -s http://localhost:8080/metrics
```

Includes events ingested, worker latency, dashboard latency, and failures.

---

## How to get a student towards mastery

**Why is `mastery` null?** The mastery endpoint returns data from the `student_mastery` table. That table is **only updated when the worker processes a `SUBMISSION_GRADED` event** for that student. So you see `"mastery": null` until at least one graded submission (with a **score** and **standard_ids**) has been ingested and processed.

**What counts as “mastery”:** For each `SUBMISSION_GRADED` event, the worker takes the event’s **score** (0–100), converts it to a 0–1 mastery value (`score / 100`), and upserts one row per **standard_id**. So mastery is “latest score per standard” (later graded events overwrite earlier ones for that standard).

**Steps to get non-null mastery for a student:**

1. **Worker must be running** — Mastery is written by the worker when it processes the outbox. If you use `docker compose up`, the worker is already running. If you run locally, start it with `go run ./cmd/worker` in a separate terminal.

2. **Send a `SUBMISSION_GRADED` event** for that student with both **score** and **standard_ids** (and the usual required fields). Example for `student-1`:

   ```bash
   curl -s -X POST http://localhost:8080/events \
     -H "Content-Type: application/json" \
     -d '{
       "event_id": "ev-mastery-demo-1",
       "source": "manual",
       "timestamp": "2025-02-13T14:00:00Z",
       "student_id": "student-1",
       "class_id": "class-1",
       "assignment_id": "assign-demo",
       "standard_ids": ["std-math-1", "std-math-2"],
       "score": 85.0,
       "type": "SUBMISSION_GRADED"
     }'
   ```

   Expect **202 Accepted**.

3. **Wait a few seconds** for the worker to claim the outbox row and run the mastery update.

4. **Call the mastery endpoint again:**

   ```bash
   curl -s http://localhost:8080/students/student-1/mastery | jq
   ```

   You should see something like:

   ```json
   {
     "student_id": "student-1",
     "mastery": [
       { "standard_id": "std-math-1", "mastery_score": 0.85 },
       { "standard_id": "std-math-2", "mastery_score": 0.85 }
     ]
   }
   ```

   (Score 85 → mastery 0.85 per standard.)

**If you already ran the simulator:** The simulator posts `SUBMISSION_GRADED` events for `student-1` … `student-20`. If you still see `mastery: null` for `student-1`, check that the worker is running and that it has had time to process the outbox (watch worker logs for processing). If the worker started after the simulator finished, it should pick up pending rows; if the DB was recreated without re-running the simulator, post a graded event as above.

---

## Quick verification checklist

| Step | What to do | What you should see |
|------|------------|----------------------|
| 1 | `docker compose up` or start API + worker | API on :8080, worker processing logs |
| 2 | `go run ./cmd/simulate` | "Posted … events" and "Dashboard response status: 200" (or 200 after worker catches up) |
| 3 | `curl …/teachers/teacher-1/classes/class-1/dashboard` | JSON with completion rate, avg score, at_risk_students, etc. |
| 4 | `curl …/students/student-1/mastery` | JSON with standard IDs and scores |
| 5 | `curl …/metrics` | Prometheus text with `events_ingested_total`, etc. |

---

## Troubleshooting

- **Dashboard empty or 404**  
  Run the simulator, then wait a few seconds for the worker to process the outbox. Hit the dashboard again.

- **"connection refused" to Postgres**  
  Ensure Postgres is up: `docker compose up -d postgres` and that `DATABASE_URL` points at localhost:5432 when running API/worker locally.

- **Migrations failed**  
  With Docker: `docker compose down -v` then `docker compose up` to start with a fresh DB. Locally: fix `DATABASE_URL` and run `migrate … up` again.

- **Simulate fails**  
  Ensure the API is running on the same host/port (default http://localhost:8080). Set `API_URL` if the API is elsewhere.

---

## Next steps

- Change at-risk rules in `internal/risk`.
- Add or tweak event types and dashboard fields in `internal/domain` and handlers.
- Run the simulator with higher `STUDENTS`/`ASSIGNMENTS` for load (e.g. `STUDENTS=100 ASSIGNMENTS=50 go run ./cmd/simulate`).
- Point Prometheus at `http://localhost:8080/metrics` and build dashboards.
