# LVBP Project - Backend

This directory contains the Go backend for the LVBP project, built using Hexagonal Architecture.

## Prerequisites
- [Go](https://go.dev/doc/install) 1.22 or higher
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) for type-safe SQL code generation

## Tech Stack
- **Language**: Go
- **Architecture**: Hexagonal (Ports & Adapters)
- **Database**: PostgreSQL 16 (via `pgx/v5` and `sqlc`)
- **Event Bus / Cache**: Redis 7 (via `go-redis/v9`)
- **Testing**: Testcontainers-Go (for real DB/Redis integration tests)

## Local Development Setup

We use `docker-compose` to spin up a fully configured local database and Redis instance. The local setup automatically injects the schema and a pseudo-dataset (mock data) so you can start developing immediately.

1. **Start the Infrastructure (DB & Redis):**
   Run the following command from the *root of the repository* (not the backend folder):
   ```bash
   docker-compose up -d
   ```
   > **Note:** The `docker-compose.override.yml` ensures that the local database is automatically populated with LVBP teams, a mock scheduled game, and an in-progress game using `seed.sql`.

2. **Download Go Dependencies:**
   Navigate into the `backend/` directory and install the required modules:
   ```bash
   cd backend
   go mod tidy
   ```

3. **Generate Database Code (sqlc):**
   If you modify the schema or any SQL queries in `backend/sql/queries.sql`, you must regenerate the Go database adapters:
   ```bash
   sqlc generate
   ```

4. **Run the Backend Server:**
   ```bash
   go run cmd/api/main.go
   ```

## Running Tests

We use `testcontainers-go` to spin up ephemeral Docker containers during testing to ensure real integration. You must have the Docker daemon running to execute these tests.

```bash
# Run all tests (including integration tests)
go test ./... -v
```

## Production Deployment

The production deployment uses the `Dockerfile` inside this directory to compile a minimal Alpine-based container. The `docker-compose.yml` file is ready for production as it mounts the schema without the pseudo-data seed.

To build the image manually:
```bash
docker build -t lvbp-backend .
```
