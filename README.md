# Comprehensive Go Tutorial

A hands-on tutorial project for learning Go (Golang), designed for developers with C/C++, JavaScript, and TypeScript experience.

**This repo also includes a production-style backend:** the **[Student Progress & Mastery Service](edtech-mastery-service/README.md)** — an event-driven Go API for learning events, student mastery, and teacher dashboards. See [Standalone Project: EdTech Mastery Service](#standalone-project-edtech-mastery-service) below.

---

## Prerequisites

- **Go 1.21+** — [Install Go](https://go.dev/doc/install)
- Verify: `go version`

## How to Use This Tutorial

Each lesson is a self-contained, runnable program. Work through them in order:

| Lesson | Topic | Run Command |
|--------|-------|-------------|
| 01 | Basics: Hello World, variables, types | `go run 01-basics/main.go` |
| 02 | Control flow: if, for, switch | `go run 02-control-flow/main.go` |
| 03 | Functions: multiple returns, variadic | `go run 03-functions/main.go` |
| 04 | Pointers & structs | `go run 04-pointers-structs/main.go` |
| 05 | Interfaces | `go run 05-interfaces/main.go` |
| 06 | Collections: slices & maps | `go run 06-collections/main.go` |
| 07 | Error handling | `go run 07-error-handling/main.go` |
| 08 | Concurrency: goroutines & channels | `go run 08-concurrency/main.go` |
| 09 | Packages & modules | `go run 09-packages/main.go` |
| 10 | Standard library | `go run 10-standard-library/main.go` |
| 11 | Generics (bonus) | `go run 11-generics/main.go` |
| project | CLI tool — culminating project | `go run project/main.go` |

## Standalone Project: EdTech Mastery Service

In addition to the tutorial lessons, this repo contains a **separate, production-style backend** in its own Go module:

| Project | Description | How to Run |
|---------|-------------|------------|
| **edtech-mastery-service** | Event-driven Go backend: learning event ingestion, student mastery & class rollups, teacher dashboards. Uses Postgres, Redis, chi, zerolog, Prometheus. | See [edtech-mastery-service/README.md](edtech-mastery-service/README.md). Quick start: `cd edtech-mastery-service && docker compose up` |

It demonstrates idiomatic Go, event-driven design, worker pools, idempotent APIs, and observability. It has its own `go.mod`, so run and build from inside `edtech-mastery-service/`.

## Quick Reference: Go vs Your Background

| Concept | C/C++ | JavaScript/TypeScript | Go |
|---------|-------|------------------------|-----|
| Memory | Manual (new/delete) | Garbage collected | Garbage collected |
| Types | Static, explicit | Dynamic / Optional | Static, explicit (with inference) |
| Pointers | Full control, complex | None (references) | Explicit but simpler |
| Concurrency | Threads, mutexes | Single-threaded + async | Goroutines, channels |
| Error handling | Exceptions / return codes | try/catch, throw | Explicit return values |
| Inheritance | Classes, virtual | Classes, prototype | Composition + interfaces |

## Project Structure

```
GoLangProject/
├── README.md                 # This file
├── go.mod                    # Module definition (tutorial + project/)
├── 01-basics/                # Variables, types, I/O
├── 02-control-flow/          # Conditionals, loops
├── 03-functions/             # Functions, multiple returns
├── 04-pointers-structs/      # Pointers and structs
├── 05-interfaces/            # Interface-based design
├── 06-collections/           # Slices, maps
├── 07-error-handling/        # Error handling patterns
├── 08-concurrency/           # Goroutines, channels
├── 09-packages/              # Multi-package structure
├── 10-standard-library/     # JSON, HTTP, files, time
├── 11-generics/              # Generics (Go 1.18+)
├── project/                  # Final CLI project
└── edtech-mastery-service/   # Standalone backend (own go.mod, Docker, APIs)
```

## Tips for Your Background

- **From C/C++**: Go has no `malloc`/`free` — the runtime handles it. Struct embedding replaces inheritance. Use `interface{}` sparingly; generics (Go 1.18+) are preferred.
- **From JS/TS**: Go is compiled and strictly typed. No `undefined`/`null` — use zero values and the comma-ok idiom. Concurrency is built-in (goroutines), not bolted on.

## Running All Lessons

```bash
# Run a specific lesson
go run 01-basics/main.go

# Run tests (where available)
go test ./...
```
