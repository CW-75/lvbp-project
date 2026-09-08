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

## 5. Agent Implementation Rules

### 5.1. Backend Constraints (Go / Postgres / Redis)
- **Strict DB Interaction:** Use `sqlc` for all queries. No raw `database/sql` queries or ORMs like Gorm are allowed.
- **SSE Handling:** For the Real-Time Channel, utilize `http.Flusher` natively via `go-chi`. You MUST enforce the `X-Accel-Buffering: no` and `Cache-Control: no-cache` headers to comply with NFR-3.
- **Isolated Business Logic:** The State Engine must not infer physical outcomes (e.g., deducing a strike by coordinates). Strictly follow the Scorekeeper's explicit input as defined in the PRD.

### 5.2. Frontend Constraints (Next.js / React)
- **Strict State Segregation:** Use `Zustand` EXCLUSIVELY for live SSE states (in-progress games). Use `TanStack Query` EXCLUSIVELY for fetching historical/static REST endpoints. Never mix these responsibilities.
- **App Router Paradigms:** Default to React Server Components. The `"use client"` directive must only be used in component trees requiring pure interactivity or consuming Zustand/TanStack hooks (e.g., Scorekeeper's visual matrix).
- **Styling Restrictions:** Custom CSS files are prohibited (except initialization). Use `Tailwind CSS` utility classes and `shadcn/ui` components for all styling, including complex SVG manipulations for the diamond and strike zone.
- **Tailwind CSS Clean Syntax & Design Tokens:**
  - **Canonical Token Utilities:** Always use standard Tailwind theme utility classes (e.g., `text-muted-foreground`, `bg-surface`, `text-primary`, `bg-background`, `border-border`, `text-foreground`) instead of arbitrary CSS variable syntax (e.g., `text-[var(--muted-foreground)]`, `bg-[var(--surface)]`, `text-[var(--primary)]`).
  - **Forbidden Redundancies:** Arbitrary bracket notation (`[var(--...)]` or `[color:...]`) is strictly prohibited when the corresponding design token or utility is declared in `@theme` / `globals.css`.
  - **Class Composition:** Use `cn()` (`clsx` + `tailwind-merge`) for dynamic and conditional class combinations to prevent style collisions and keep markup clean.
