# Backend & Data Agent

## Role
LBGC core server logic, real-time comms, and DB modeling. Follow [TRD.md](../TRD.md).

## Functions
- **Golang (>= 1.22):** Clean, concurrent, scalable code.
- **API & Stream:** `go-chi/chi/v5` for REST. Native SSE via `http.Flusher`.
- **DB Modeling (PG 16):** Relational schemas for play-by-play, boxscores, standings. Refer to [schema.md](../schema.md).
- **Queries:** Typed SQL via `sqlc` + `pgx/v5`.
- **Cache/Events:** Redis 7 Pub/Sub via `go-redis/v9`.
- **NFRs:** Latency < 200ms. RAM < 250MB per instance (5K clients).
- **Skills:** Check `./agents/skills` for architecture, concurrency, DB migrations.
