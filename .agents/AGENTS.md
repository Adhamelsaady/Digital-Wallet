# Digital Wallet & Double-Entry Ledger Engine — Agent Rules

## Project Overview

This is a **learning-focused Go project** building a production-grade Digital Wallet with a
Double-Entry Ledger Engine backed by PostgreSQL. The primary goal is to learn Go idioms,
patterns, and best practices through building a real-world financial system.

---

## Tech Stack

| Concern | Choice | Notes |
|---|---|---|
| Language | **Go** | Learning goal; excellent fit for financial services |
| Database | **PostgreSQL** | Row locking, `SERIALIZABLE` isolation, exact numeric types |
| DB Driver | **`pgx`** (native interface preferred, or via `database/sql`) | Fast, well-maintained, no ORM magic |
| Migrations | **`golang-migrate`** | SQL-file-based, no framework lock-in |
| HTTP Router | **`chi`** or plain `net/http` (Go 1.22+ routing) | Keep it minimal; ~10 endpoints don't need a big framework |
| Config | **Environment variables** via `os.Getenv` + `config.go`, or `viper` for more structure | |
| Testing | **Standard `testing` package** + `testify` for assertions | |
| Background Jobs | **In-process goroutine + ticker** to start; upgrade to `river` or `asynq` only if needed | Don't over-engineer the MVP |
| UUIDs | **`google/uuid`** | |
| Money Type | **`int64` (minor units / cents)** — no decimal library needed | Never use `float64` for money |

---

## Project Conventions

### Go Style
- Follow standard Go conventions: `gofmt`, `golint`, `go vet`.
- Use idiomatic Go: short variable names in small scopes, longer names for package-level identifiers.
- Prefer explicit error handling — never silently ignore errors.
- Use `errors.Is` / `errors.As` for error inspection; define sentinel errors or typed errors per package.
- Keep functions small and focused; extract helpers rather than growing large functions.
- Use Go's standard `context.Context` for cancellation and deadline propagation through every DB call and HTTP handler.

### Project Structure (target layout)
```
.
├── cmd/
│   └── server/          # main entrypoint
├── internal/
│   ├── config/          # env-based config loading
│   ├── db/              # pgx connection pool setup
│   ├── ledger/          # double-entry ledger domain logic
│   ├── wallet/          # wallet domain logic
│   ├── api/             # HTTP handlers and routing (chi)
│   └── worker/          # background goroutines (webhooks, reconciliation)
├── migrations/          # SQL migration files (golang-migrate format)
├── .agents/             # agent customizations (this file)
├── go.mod
└── go.sum
```

### Database
- All money values are stored as `BIGINT` (minor units, e.g., cents). **Never use `FLOAT` or `NUMERIC` with decimals for balances.**
- Every financial mutation must be a **double-entry** transaction: debits equal credits.
- Use PostgreSQL **`SERIALIZABLE`** isolation or explicit `SELECT ... FOR UPDATE` row locking for balance updates to prevent race conditions.
- Migration files follow the `golang-migrate` naming convention: `000001_create_wallets.up.sql` / `000001_create_wallets.down.sql`.
- Never use an ORM. Write raw SQL queries. Use `pgx` named or positional parameters.

### API Design
- RESTful endpoints; JSON request/response bodies.
- Return meaningful HTTP status codes (201 for creation, 422 for validation errors, 409 for conflicts, etc.).
- All responses include a consistent envelope structure where practical.
- Validate inputs at the handler layer before touching the database.

### Testing
- Write unit tests alongside every non-trivial function (`_test.go` in the same package).
- Use `testify/assert` and `testify/require` for assertions.
- Integration tests that need Postgres should use a real test database or `testcontainers-go`.
- Table-driven tests are preferred for functions with multiple input/output cases.

### Error Handling
- Domain errors (e.g., `ErrInsufficientFunds`, `ErrWalletNotFound`) must be defined as typed or sentinel errors.
- HTTP handlers map domain errors to appropriate HTTP status codes.
- Never leak internal DB errors to API responses.

### Background Jobs
- Start with `time.Ticker` + goroutines for webhooks and reconciliation tasks.
- All goroutines must respect `context.Context` for graceful shutdown.
- Upgrade to a proper queue (`river`, `asynq`) only when complexity demands it.

---

## Learning Focus

Since this project is primarily for **learning Go**, agents should:

1. **Explain non-obvious Go patterns** when introducing them (e.g., embedding, interface satisfaction, table-driven tests).
2. **Prefer idiomatic Go** over patterns from other languages — avoid Java/Python-style abstractions.
3. **Keep the code readable** over being overly clever. Clarity > brevity when in conflict.
4. **Don't over-engineer**: avoid generics, reflection, or complex abstractions unless there is a clear, demonstrated need.
5. **Introduce concepts progressively**: start simple, refactor to better patterns as complexity grows.
6. **Add comments** on non-obvious decisions (e.g., why `SERIALIZABLE`, why `int64` for money).

---

## Agent Behavioral Rules

- **Always use `int64` for money amounts.** Never suggest `float64` or `decimal` libraries for this project.
- **Never use an ORM.** All database interactions must use raw SQL via `pgx`.
- **Always handle errors explicitly.** Do not use `_` to discard errors from DB calls or I/O operations.
- **Always pass `context.Context` as the first argument** to any function that performs I/O.
- **Do not add dependencies** without a clear justification. This is an MVP; keep `go.mod` lean.
- **Write migrations as plain SQL files** in `migrations/`. Do not use Go-based migration DSLs.
- **Prefer table-driven tests** for any function with more than one logical case.
- When proposing new packages or modules, **explain why** relative to the learning goal.
