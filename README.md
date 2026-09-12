# Digital Wallet & Double-Entry Ledger Engine

## Overview

This project implements a robust financial system designed to handle account balances and transactions securely and accurately. The core accounting engine adheres strictly to **double-entry bookkeeping principles**, ensuring that every financial mutation balances debits against credits.

## Key Architectural Decisions

- **Double-Entry Accounting:** Every transaction consists of matched debit and credit entries.
- **Exact Integer Financials:** All monetary amounts are stored in minor units (`int64`, e.g., cents) to completely avoid floating-point rounding errors.
- **Strong Concurrency & Isolation:** Built with PostgreSQL `SERIALIZABLE` transactions and explicit row-locking to guarantee data integrity under concurrent traffic.
- **Idiomatic & No-ORM Go:** Pure SQL via `pgx` for full control over database execution and zero magic abstractions.

## Tech Stack

- **Language:** Go 1.22+
- **Database:** PostgreSQL 16
- **Database Driver / Connection Pool:** `pgx/v5` (`pgxpool`)
- **HTTP Router:** `chi/v5`
- **UUID Generation:** `google/uuid`
- **Testing:** Standard `testing` package + `testify` for assertions
- **Containerization:** Docker & Docker Compose

---

## 📁 Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point & dependency wiring
├── internal/
│   ├── config/                     # Environment configuration loader
│   ├── db/                         # pgx connection pool setup & health ping
│   ├── ledger/                     # Domain entities, validation, and business logic
│   ├── storage/                    # Database repositories with raw SQL
│   └── api/                        # HTTP handlers, router, and JSON response helpers
├── migrations/                     # SQL schema migration files
├── docker-compose.yml              # Local PostgreSQL 16 container definition
├── .env                            # Local environment variables (ignored in git)
├── go.mod
└── go.sum
```

---

## Features Implemented So Far

### 1. Skeleton & Local Environment (Phase 1)
- [x] Dockerized PostgreSQL 16 database with automated healthchecks.
- [x] Resilient database connection pooling with `pgxpool.New` and startup ping verification.
- [x] HTTP server with `chi` router, request logger, and panic recovery middleware.
- [x] Health check endpoint: `GET /health` → `{"status":"ok"}`.

### 2. Schema & Accounts (Phase 2)
- [x] Database migration for `accounts` table with unique constraint on `(owner_id, currency)` and hash indexing on `owner_id`.
- [x] Domain entity `Account` with constructor validation (ISO currency check, valid owner UUID).
- [x] Table-driven unit tests for account creation and validation logic using `testify`.
- [x] Storage repository (`AccountRepository`) with raw SQL queries and PostgreSQL error code mapping (handling `23505` unique violations and `ErrNoRows`).
- [x] API handler for creating accounts: `POST /accounts`.

---

## Getting Started

### 1. Prerequisites
- [Go](https://go.dev/dl/) (1.22 or higher)
- [Docker & Docker Compose](https://www.docker.com/)
- `curl` or Postman for API testing

### 2. Start PostgreSQL
```powershell
docker compose up -d
```

### 3. Run Database Migrations
Apply the initial database schema:
```powershell
docker exec -i digital_wallet_db psql -U wallet -d digital_wallet < migrations/000001_create_accounts_table.up.sql
```

### 4. Run Unit Tests
```powershell
go test -v ./internal/ledger/...
```

### 5. Start the API Server
```powershell
go run ./cmd/api
```

---

## 📡 API Endpoints

| Method | Path | Description | Status Code |
|---|---|---|---|
| `GET` | `/health` | Server health check | `200 OK` |
| `POST` | `/accounts` | Create a new financial account | `201 Created` |

### Example: Create an Account
```powershell
curl -X POST http://localhost:8080/accounts `
  -H "Content-Type: application/json" `
  -d '{"owner_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "currency": "USD", "type": "AVAILABLE"}'
```

**Response (`201 Created`):**
```json
{
  "id": "e8222bb5-c53b-4ef8-bb6d-6bb9bd380a11",
  "owner_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
  "currency": "USD",
  "type": "AVAILABLE",
  "created_at": "2026-09-12T01:19:10Z",
  "updated_at": "2026-09-12T01:19:10Z"
}
```