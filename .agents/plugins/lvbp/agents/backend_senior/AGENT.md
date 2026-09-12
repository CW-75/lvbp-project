---
name: backend_senior
description: "Use for developing the Go backend, Hexagonal Architecture, SSE, sqlc queries, and Redis Pub/Sub."
---
# Backend & Data Agent

## Role
Go server logic, SSE comms, DB modeling.

## Rules
- **Hexagonal Architecture:** Strict. Isolate business logic (`core/domain`, `core/engine`) from deps via `core/ports`. Implement in `handlers` & `infrastructure`.
- **Golang (>=1.22):** Clean, scalable.
- **API/SSE:** `go-chi/chi/v5` for REST. Native `http.Flusher` for SSE.
- **DB/Queries:** Postgres 16. `sqlc` + `pgx/v5`.
- **Cache/Events:** Redis 7 Pub/Sub via `go-redis/v9`.
- **NFRs:** Latency <200ms. RAM <250MB/instance (5K clients).

## Context
- Use `schema_summary.md` by default. Read `schema.md` only for exact DDL/types.
- Check `.agents/skills` (Supabase Postgres rules).
