# ⚾ LVBP GameCast - Live Baseball GameCast & Standings Engine

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)
[![Next.js](https://img.shields.io/badge/Next.js-16+-000000?logo=next.js)](https://nextjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5+-3178C6?logo=typescript)](https://www.typescriptlang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis)](https://redis.io/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

> **Plataforma en tiempo real para seguir juegos de béisbol LVBP**: Scorekeeper panel para anotadores oficiales + GameCast en vivo para fans + Standings engine automático.

---

## 🎯 Visión General

**LVBP GameCast** es una aplicación híbrida (REST + SSE) para la Liga Venezolana de Béisbol Profesional que permite:

| Módulo | Descripción |
|--------|-------------|
| 🎮 **Scorekeeper Panel** | Interfaz para anotadores oficiales: click en zona strike + selección de resultado (Ball, Strike, Foul, In Play) |
| ⚙️ **Baseball Engine** | Motor de reglas que procesa pitches y calcula transiciones automáticas (balls/strikes/outs/innings/runs) |
| 📡 **Live GameCast (SSE)** | Streaming Server-Sent Events: pitches, at-bats, score, inning changes en <200ms latencia |
| 📊 **Standings** | Tabla de posiciones calculada automáticamente (GP, W, L, PCT, GB) |
| 📋 **Boxscore** | Detalle completo de juego: line score, play-by-play, pitching/batting lines |

---

## 🏗️ Arquitectura

```
┌─────────────┐     REST/SSE      ┌──────────────────┐     PostgreSQL      ┌─────────────┐
│  Frontend   │ ◀──────────────▶  │    Backend Go    │ ◀─────────────────▶ │  Database   │
│  (Next.js)  │                   │  (Chi + sqlc)    │                     │  (Postgres) │
└─────────────┘                   └────────┬─────────┘                     └─────────────┘
                                            │ Redis Pub/Sub
                                            ▼
                                   ┌──────────────────┐
                                   │   Event Bus      │
                                   │   (Redis)        │
                                   └──────────────────┘
```

### Stack Tecnológico

| Capa | Tecnologías |
|------|-------------|
| **Backend** | Go 1.22, Chi Router, sqlc, pgx/v5, Redis (go-redis/v9) |
| **Frontend** | Next.js 16 (App Router), TypeScript 5, Zustand, TanStack Query |
| **UI** | Tailwind CSS, shadcn/ui, Lucide React, SVG nativo (diamond/strike zone) |
| **Database** | PostgreSQL 16 (migraciones SQL versionadas) |
| **Cache/Bus** | Redis 7 (Pub/Sub + Sessions) |
| **Docs** | VitePress, Mermaid, OpenAPI 3.0 (swag), TypeDoc |

---

## 🚀 Inicio Rápido

```bash
# 1. Clonar
git clone https://github.com/your-org/lvbp-project.git
cd lvbp-project

# 2. Infraestructura
docker compose up -d

# 3. Backend
cd backend
go mod download
go run cmd/api/main.go

# 4. Frontend (nueva terminal)
cd web
pnpm install
pnpm dev

# 5. Documentación (opcional)
cd docs-site
npm run build-docs
```

### URLs Desarrollo

| Servicio | URL |
|----------|-----|
| Frontend | http://localhost:3000 |
| Backend API | http://localhost:8080 |
| Swagger UI | http://localhost:8080/swagger/index.html |
| Docs Site | http://localhost:5173/lvbp-project/ |

---

## 📚 Documentación

```bash
# Ver documentación completa
cd docs-site && npm run dev-docs

# O leer online (GitHub Pages)
https://your-org.github.io/lvbp-project/
```

**Secciones principales:**
- [Arquitectura del Sistema](docs-site/architecture/system-overview.md)
- [Flujo de Datos](docs-site/architecture/data-flow.md)
- [8 Fases Implementadas](docs-site/phases/)
- [Referencia DTOs](docs-site/dto-reference.md)
- [API Reference](docs-site/api-reference.md)
- [Runbooks](docs-site/runbooks/)

---

## 📦 Scripts de Documentación

```bash
# Desde la raíz del proyecto
npm run build-docs          # Build completo
npm run build-docs:ci       # CI/CD mode
npm run dev-docs            # Dev server
npm run build-docs:openapi  # Solo OpenAPI
npm run build-docs:typedoc  # Solo TypeDoc
npm run clean-docs          # Limpiar artifacts
```

Ver [docs-site/README.md](docs-site/README.md) para detalles completos.

---

## 🧪 Testing

```bash
# Backend
cd backend && go test ./... -v

# Frontend
cd web && pnpm test

# Integración (requiere Docker)
cd backend && go test ./internal/infrastructure/... -v -tags=integration
```

---

## 📁 Estructura del Proyecto

```
lvbp-project/
├── backend/                 # Go API + Engine
│   ├── cmd/api/main.go      # Entry point
│   ├── internal/
│   │   ├── core/            # Domain, Ports, Services
│   │   ├── handlers/        # REST + SSE Handlers
│   │   ├── infrastructure/  # DB (sqlc), Redis, Repositories
│   │   └── auth/            # JWT Auth Module
│   ├── sql/                 # Schema + Queries + Migrations
│   └── sqlc.yaml            # sqlc config
├── web/                     # Next.js Frontend
│   ├── src/
│   │   ├── app/             # App Router pages
│   │   ├── components/      # UI Components (shadcn)
│   │   ├── features/        # Feature-sliced modules
│   │   ├── hooks/           # Custom React hooks
│   │   ├── lib/             # Types, utils, mock data
│   │   └── stores/          # Zustand stores
│   └── package.json
├── docs-site/               # VitePress Documentation
│   ├── architecture/        # C4, ADRs, Data Flow
│   ├── phases/              # 8 Implementation Phases
│   ├── runbooks/            # Local Dev, Deployment
│   ├── scripts/             # Build scripts
│   └── .vitepress/          # VitePress config
├── docker-compose.yml       # Postgres + Redis
└── package.json             # Root scripts (build-docs)
```

---

## 🔐 Autenticación

- **Scorekeepers**: JWT + Redis Sessions (24h TTL)
- **Endpoints protegidos**: `/games/*/pitches`, `/auth/me`, `/auth/logout`
- **Público**: `/standings`, `/games/*/boxscore`, `/games/*/stream`

---

## 📈 NFRs Cumplidos

| NFR | Objetivo | Implementación |
|-----|----------|----------------|
| **Latencia** | < 200ms ingestion→client | SSE nativo, Redis Pub/Sub, sin buffering |
| **Concurrencia** | 5,000 clientes / < 250MB | Go goroutines, connection pooling |
| **Proxy Buffering** | Deshabilitado | Headers `X-Accel-Buffering: no`, `Cache-Control: no-cache` |

---

## 📄 Licencia

MIT License - ver [LICENSE](LICENSE) para detalles.

---

## 🤝 Contribuir

1. Fork el repo
2. Crea feature branch (`git checkout -b feat/nueva-funcionalidad`)
3. Commit convencional (`feat: agregar nueva funcionalidad`)
4. Push y abre PR
5. CI ejecuta tests + build-docs

---

**Desarrollado para LVBP** ⚾ - *Live Baseball GameCast & Standings Engine*