# Runbook: Desarrollo Local

## Prerrequisitos

| Herramienta | Versión | Instalación |
|-------------|---------|-------------|
| Go | 1.22+ | `winget install GoLang.Go` / `brew install go` |
| Node.js | 20+ | `winget install OpenJS.NodeJS` / `brew install node` |
| pnpm | 9+ | `npm install -g pnpm` |
| Docker | 24+ | Docker Desktop |
| sqlc | 1.27+ | `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest` |
| swag | latest | `go install github.com/swaggo/swag/cmd/swag@latest` |
| Make | - | Incluido en Docker Desktop / `choco install make` |

## Inicio Rápido

```bash
# 1. Clonar y entrar
cd lvbp-project

# 2. Levantar infraestructura (PostgreSQL + Redis)
docker compose up -d

# 3. Backend: instalar deps y generar código
cd backend
go mod download
sqlc generate
swag init -g cmd/api/main.go -o ../docs-site/public/api --parseDependency --parseInternal

# 4. Frontend: instalar deps
cd ../web
pnpm install

# 5. Ejecutar todo (terminales separadas)
# Terminal 1 - Backend
cd backend && go run cmd/api/main.go

# Terminal 2 - Frontend
cd web && pnpm dev

# Terminal 3 - Docs (opcional)
cd docs-site && pnpm install && pnpm dev
```

## URLs de Desarrollo

| Servicio | URL | Descripción |
|----------|-----|-------------|
| Frontend | http://localhost:3000 | Next.js App |
| Backend API | http://localhost:8080 | REST + SSE |
| API Docs (Swagger) | http://localhost:8080/swagger/index.html | Swagger UI |
| Docs Site | http://localhost:5173 | VitePress |
| PostgreSQL | localhost:5432 | user: lvbp, pass: lvbp, db: lvbp |
| Redis | localhost:6379 | sin auth |

## Variables de Entorno

### Backend (`.env` en `backend/` o root)
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=lvbp
DB_PASSWORD=lvbp
DB_NAME=lvbp
DB_SSLMODE=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Auth
JWT_SECRET=your-super-secret-key-change-in-production
JWT_EXPIRATION=24h

# Server
SERVER_PORT=8080
GIN_MODE=debug
```

### Frontend (`.env.local` en `web/`)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080
```

## Comandos Make Útiles

```bash
# Raíz del proyecto
make help                    # Ver todos los comandos
make dev                     # Levantar todo (docker + backend + frontend)
make test                    # Tests backend + frontend
make lint                    # Lint todo
make sqlc-generate           # Regenerar sqlc
make swag-generate           # Regenerar OpenAPI
make db-migrate              # Ejecutar migraciones
make db-seed                 # Seed data inicial
make clean                   # Limpiar build artifacts
```

## Estructura de Tests

### Backend
```bash
cd backend

# Todos los tests
go test ./... -v

# Solo unit tests (rápidos)
go test ./internal/core/... ./internal/auth/... -v

# Tests de integración (requiere Docker)
go test ./internal/infrastructure/... -v -tags=integration

# Coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Frontend
```bash
cd web

# Unit tests
pnpm test

# E2E tests (Playwright)
pnpm test:e2e

# Type check
pnpm type-check

# Lint
pnpm lint
```

## Flujo de Trabajo Git

```bash
# 1. Crear feature branch
git checkout -b feat/nueva-funcionalidad

# 2. Desarrollar + tests
# ...

# 3. Pre-commit hooks (auto)
# - go fmt / go vet
# - pnpm lint --fix
# - conventional commit message

# 4. Commit
git add .
git commit -m "feat: agregar nueva funcionalidad"

# 5. Push y PR
git push origin feat/nueva-funcionalidad
# Crear PR en GitHub
```

## Migraciones de Base de Datos

```bash
# Crear nueva migración
cd backend
migrate create -ext sql -dir sql/migrations -seq nombre_migracion

# Aplicar migraciones
migrate -path sql/migrations -database "postgres://lvbp:lvbp@localhost:5432/lvbp?sslmode=disable" up

# Rollback última
migrate -path sql/migrations -database "postgres://lvbp:lvbp@localhost:5432/lvbp?sslmode=disable" down 1
```

## Regenerar Código Generado

```bash
# sqlc (después de cambiar schema.sql o queries.sql)
cd backend && sqlc generate

# swag/OpenAPI (después de cambiar annotations en handlers)
cd backend && swag init -g cmd/api/main.go -o ../docs-site/public/api --parseDependency --parseInternal

# typedoc (después de cambiar types TS)
cd web && npx typedoc --out ../docs-site/public/typedoc --entryPoints src/lib
```

## Troubleshooting Común

### Backend no conecta a PostgreSQL
```bash
# Verificar contenedor
docker compose ps
docker compose logs postgres

# Verificar conexión manual
psql -h localhost -p 5432 -U lvbp -d lvbp
```

### Puerto 8080 ocupado
```bash
# Matar proceso
lsof -ti:8080 | xargs kill -9
# O cambiar puerto en .env
```

### sqlc generate falla
```bash
# Verificar schema.sql syntax
# Verificar queries.sql referencia tablas/columnas correctas
sqlc vet
```

### Frontend build falla
```bash
cd web
rm -rf node_modules .next
pnpm install
pnpm build
```

### Tests de integración fallan
```bash
# Asegurar Docker corriendo
docker compose up -d
# Verificar testcontainers puede crear contenedores
docker run --rm hello-world
```

## Debugging

### Backend (VS Code / GoLand)
```json
// .vscode/launch.json
{
  "configurations": [
    {
      "name": "Launch Backend",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/backend/cmd/api/main.go",
      "env": {
        "DB_HOST": "localhost",
        "DB_PORT": "5432",
        "DB_USER": "lvbp",
        "DB_PASSWORD": "lvbp",
        "DB_NAME": "lvbp",
        "REDIS_ADDR": "localhost:6379"
      }
    }
  ]
}
```

### Frontend (VS Code)
```json
// .vscode/launch.json (añadir)
{
  "name": "Next.js: debug",
  "type": "node",
  "request": "launch",
  "program": "${workspaceFolder}/web/node_modules/next/dist/bin/next",
  "args": ["dev"],
  "cwd": "${workspaceFolder}/web"
}
```

## Logs Estructurados

Backend usa `slog` (Go 1.21+):
```go
slog.Info("pitch_processed", 
    "game_id", gameID, 
    "pitch_result", pitch.PitchResult,
    "latency_ms", elapsed.Milliseconds())
```

Ver logs en terminal donde corre `go run` o:
```bash
# Si corriendo en Docker
docker compose logs -f api
```

## Performance Profiling

```bash
# CPU Profile
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

# Memory Profile
go tool pprof http://localhost:8080/debug/pprof/heap

# Goroutines
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

## Seed Data Desarrollo

```bash
# Ejecutar seed.sql
cd backend
psql -h localhost -U lvbp -d lvbp -f sql/seed.sql
```

Crea:
- 8 equipos LVBP
- Usuario scorekeeper de prueba: `test@lvbp.com` / `password123`
- Juegos de ejemplo (scheduled, in_progress, final)