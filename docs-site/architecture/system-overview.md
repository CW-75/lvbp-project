# Visión General de la Arquitectura

## Diagrama de Contexto (C4 Level 1)

```mermaid
C4Context
title Contexto del Sistema LVBP GameCast

Person(scorekeeper, "Scorekeeper", "Usuario autenticado que registra jugadas en tiempo real desde el estadio")
Person(fan, "Fan/Espectador", "Consume scoreboard en vivo y standings desde cualquier dispositivo")

System_Boundary(boundary, "LVBP GameCast Platform") {
  System(api, "API Gateway / Backend", "Go + Chi + PostgreSQL + Redis")
  System(frontend, "Frontend App", "Next.js 16 + TypeScript + Tailwind")
}

System_Ext(db, "PostgreSQL 16", "Persistencia: Games, AtBats, Pitches, Teams, Users")
System_Ext(redis, "Redis 7", "Cache + Event Bus Pub/Sub + Sessions")
System_Ext(auth, "Auth Provider", "JWT + Redis Sessions (interno)")

Rel(scorekeeper, frontend, "HTTPS - Scorekeeper Panel")
Rel(fan, frontend, "HTTPS - GameCast / Standings")
Rel(frontend, api, "REST + SSE")
Rel(api, db, "SQL (sqlc/pgx)")
Rel(api, redis, "Pub/Sub + Sessions")
Rel(api, auth, "Validación JWT")
```

## Diagrama de Contenedores (C4 Level 2)

```mermaid
C4Container
title Contenedores - LVBP Backend

Container(api, "API Server", "Go 1.22 + Chi", "Expone REST + SSE endpoints")
Container(engine, "Scorekeeper Engine", "Go", "Motor de reglas de béisbol (State Machine)")
Container(repo, "Repository Layer", "Go + sqlc", "Acceso a datos tipado PostgreSQL")
Container(eventbus, "Event Bus", "Redis Pub/Sub", "Distribución eventos tiempo real")
Container(authSvc, "Auth Service", "Go", "JWT + Session Management")

ContainerDb(db, "PostgreSQL", "SQL", "Games, AtBats, Pitches, Teams, Users")
ContainerDb(redis, "Redis", "KV + Pub/Sub", "Sessions + Event Channels + Cache")

Rel(api, engine, "Usa", "ScorekeeperService interface")
Rel(api, repo, "Usa", "GameRepository interface")
Rel(api, eventbus, "Publica", "EventPublisher interface")
Rel(api, authSvc, "Valida", "Middleware JWT")
Rel(engine, repo, "Persiste", "GameRepository")
Rel(engine, eventbus, "Emite", "EventPublisher")
Rel(repo, db, "SQL", "Queries tipadas")
Rel(authSvc, redis, "Sesiones", "Redis Session Store")
Rel(authSvc, db, "Usuarios", "UserRepository")
```

## Canales de Comunicación

| Canal | Protocolo | Patrón | Uso |
|-------|-----------|--------|-----|
| **Histórico** | REST/HTTPS | Request-Response | Standings, Boxscore, Auth |
| **Tiempo Real** | SSE (Server-Sent Events) | Pub/Sub Unidireccional | Pitches, AtBats, Score, Inning Changes |
| **Interno** | Redis Pub/Sub | Fan-out | Engine → SSE Handlers |

## Principios de Diseño Aplicados

### SOLID en Backend
- **S** - Cada handler tiene responsabilidad única (Ingestion, Standings, Boxscore, SSE)
- **O** - Interfaces `ports.GameRepository`, `ports.ScorekeeperService`, `ports.EventBus` abiertas a extensión
- **L** - Implementaciones intercambiables (PostgresRepo, MockRepo para tests)
- **I** - Interfaces segregadas por capacidad (GameRepository vs StandingsRepository)
- **D** - Inyección de dependencias via constructor (`NewBoxscoreHandler(repo)`)

### KISS / YAGNI
- No ORM, solo `sqlc` para queries tipadas
- Sin framework SSE externo, `http.Flusher` nativo
- Motor de reglas sigue input explícito del scorekeeper (no infiere de coordenadas)

### Hexagonal Traits (Modular Monolith)
```
internal/
├── auth/           # Módulo autenticación
├── core/           # Dominio puro (domain, ports, services)
├── handlers/       # Adaptadores entrada (REST, SSE)
├── infrastructure/ # Adaptadores salida (Postgres, Redis)
└── pkg/            # Shared (Redis client)
```

## Decisiones Arquitectónicas Clave (ADRs)

| ADR | Título | Estado |
|-----|--------|--------|
| [001](./adr/001-hybrid-rest-sse.md) | Arquitectura Híbrida REST + SSE | Aceptado |
| [002](./adr/002-sqlc-over-orm.md) | sqlc en lugar de ORM | Aceptado |
| [003](./adr/003-scorekeeper-explicit-input.md) | Scorekeeper sigue input explícito | Aceptado |
| [004](./adr/004-redis-pubsub-eventbus.md) | Redis Pub/Sub como Event Bus | Aceptado |
| [005](./adr/005-zustand-tanstack-segregation.md) | Separación Zustand (live) / TanStack (REST) | Aceptado |

## NFRs Cumplidos

| NFR | Objetivo | Implementación |
|-----|----------|----------------|
| **NFR-1 Latencia** | < 200ms ingestion→client | SSE nativo, Redis Pub/Sub, sin buffering proxy |
| **NFR-2 Concurrencia** | 5,000 clientes / < 250MB | Go goroutines, connection pooling, streaming |
| **NFR-3 Proxy Buffering** | Deshabilitado | Headers `X-Accel-Buffering: no`, `Cache-Control: no-cache` |