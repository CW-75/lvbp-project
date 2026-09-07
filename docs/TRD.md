# Technical Requirements Document (TRD)

**Project:** Live Baseball GameCast & Standings Engine

## 1. General Architecture

Decoupled hybrid architecture:

- **Historical Channel (REST):** Static queries (Standings, Boxscore).
- **Real-Time Channel (SSE):** Unidirectional flow for pitches and plays.

## 2. Technology Stack

- **Backend:** Go (Golang) >= 1.22.
- **Routing and Streaming:** `go-chi/chi/v5` (native control of `http.Flusher`).
- **Main Database:** PostgreSQL 16 (Typed queries with `sqlc` + `pgx/v5`). Refer to [schema.md](./schema.md).
- **Cache and Event Bus:** Redis 7 (Pub/Sub with `go-redis/v9`).
- **Frontend (Framework):** Next.js 16+ (App Router, TypeScript 5+).
- **State Management (Frontend):** Zustand for SSE events (Live) and TanStack Query for REST.
- **User Interface:** Tailwind CSS + shadcn/ui + lucide-react. (Diamond and strike zone rendering via native SVG).

## 3. Non-Functional Requirements (NFR)

- **NFR-1. Latency:** < 200 ms from the scorekeeper's ingestion to the client via SSE.
- **NFR-2. Concurrency:** Support at least 5,000 concurrent clients per instance with a consumption of < 250 MB RAM.
- **NFR-3. Proxy and Buffering:** Disable intermediate buffering by sending `X-Accel-Buffering: no` and `Cache-Control: no-cache` headers.

## 4. Software Design Principles

- **SOLID:** Single responsibility per component/package, open-closed interfaces, Liskov substitution, interface segregation, dependency inversion via Go interfaces & React hooks.
- **KISS:** Keep implementations simple and direct. Avoid unnecessary abstraction layers or over-engineering.
- **YAGNI:** Implement only features required by current [PRD.md](./PRD.md)/[TRD.md](./TRD.md). No speculative abstractions or unrequested features.
- **DRY:** Share reusable logic via domain packages (Backend) and custom hooks/components (Frontend). Avoid duplicate business rules.
