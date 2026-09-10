# Digital Wallet & Double-Entry Ledger Engine

## Overview

This project implements a robust financial system designed to handle account balances and transactions securely and accurately. The core accounting engine adheres strictly to **double-entry bookkeeping principles**, ensuring that every financial mutation balances debits against credits.

## Key Architectural Decisions

- **Double-Entry Accounting:** Every transaction consists of matched debit and credit entries.
- **Exact Integer Financials:** All monetary amounts are stored in minor units (`int64`, e.g., cents) to completely avoid floating-point rounding errors.
- **Strong Concurrency & Isolation:** Built with PostgreSQL `SERIALIZABLE` transactions and explicit row-locking to guarantee data integrity under concurrent traffic.
- **Idiomatic & No-ORM Go:** Pure SQL via `pgx` for full control over database execution and zero magic abstractions.

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL 16
- **Database Driver:** `pgx` (v5)
- **Router:** `chi` (v5)