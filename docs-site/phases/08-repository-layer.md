# Fase 8: Capa Repositorio (Data Access)

## Resumen

Implementación de acceso a datos usando `sqlc` para queries type-safe PostgreSQL. Patrón Repository separa dominio de infraestructura.

## Arquitectura

```
internal/
├── core/ports/repositories.go     # Interfaces (Domain Layer)
├── infrastructure/db/
│   ├── models.go                  # sqlc generated structs
│   ├── queries.sql.go             # sqlc generated methods
│   └── db.go                      # Pool config
└── infrastructure/repositories/
    ├── game_repository.go         # Implementación GameRepository
    └── game_repository_test.go    # Tests con testcontainers
```

## Interfaces de Dominio (`core/ports/repositories.go`)

```go
type GameRepository interface {
    GetGameByID(ctx context.Context, id uuid.UUID) (*domain.Game, error)
    CreateGame(ctx context.Context, game *domain.Game) (*domain.Game, error)
    UpdateGame(ctx context.Context, game *domain.Game) (*domain.Game, error)
    GetCurrentAtBat(ctx context.Context, gameID uuid.UUID) (*domain.AtBat, error)
    CreateAtBat(ctx context.Context, atBat *domain.AtBat) (*domain.AtBat, error)
    UpdateAtBatResult(ctx context.Context, atBat *domain.AtBat) (*domain.AtBat, error)
    GetPitchesForAtBat(ctx context.Context, atBatID uuid.UUID) ([]*domain.Pitch, error)
    CreatePitch(ctx context.Context, pitch *domain.Pitch) (*domain.Pitch, error)
    GetOutsForInning(ctx context.Context, gameID uuid.UUID, inning int, isTopInning bool) (int, error)
    GetStandings(ctx context.Context) ([]*domain.TeamStandings, error)
}
```

## sqlc Configuration (`backend/sqlc.yaml`)

```yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "sql/schema.sql"
    queries: "sql/queries.sql"
    gen:
      go:
        package: "db"
        out: "internal/infrastructure/db"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_exact_table_names: false
        overrides:
          - db_type: "uuid"
            go_type: "github.com/google/uuid.UUID"
          - db_type: "timestamptz"
            go_type: "time.Time"
```

## Queries SQL (`backend/sql/queries.sql`)

```sql
-- name: GetGameByID :one
SELECT id, home_team_id, away_team_id, status, home_score, away_score,
       current_inning, is_top_inning, start_time, created_at
FROM games WHERE id = $1;

-- name: CreateGame :one
INSERT INTO games (home_team_id, away_team_id, status, home_score, away_score,
                   current_inning, is_top_inning, start_time)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, home_team_id, away_team_id, status, home_score, away_score,
          current_inning, is_top_inning, start_time, created_at;

-- name: UpdateGame :one
UPDATE games SET home_score = $2, away_score = $3, current_inning = $4,
       is_top_inning = $5, status = $6
WHERE id = $1
RETURNING id, home_team_id, away_team_id, status, home_score, away_score,
          current_inning, is_top_inning, start_time, created_at;

-- name: GetCurrentAtBat :one
SELECT id, game_id, inning, is_top_inning, batter_id, pitcher_id,
       result, runs_scored, outs_recorded, created_at
FROM at_bats
WHERE game_id = $1 AND result IS NULL
ORDER BY created_at DESC LIMIT 1;

-- name: CreateAtBat :one
INSERT INTO at_bats (game_id, inning, is_top_inning, batter_id, pitcher_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, game_id, inning, is_top_inning, batter_id, pitcher_id,
          result, runs_scored, outs_recorded, created_at;

-- name: UpdateAtBatResult :one
UPDATE at_bats SET result = $2, runs_scored = $3, outs_recorded = $4
WHERE id = $1
RETURNING id, game_id, inning, is_top_inning, batter_id, pitcher_id,
          result, runs_scored, outs_recorded, created_at;

-- name: GetPitchesForAtBat :many
SELECT id, at_bat_id, pitch_number, coordinate_x, coordinate_y,
       pitch_result, balls_before, strikes_before, outs_before,
       velocity_mph, created_at
FROM pitches
WHERE at_bat_id = $1
ORDER BY pitch_number;

-- name: CreatePitch :one
INSERT INTO pitches (at_bat_id, pitch_number, coordinate_x, coordinate_y,
                     pitch_result, balls_before, strikes_before, outs_before,
                     velocity_mph)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, at_bat_id, pitch_number, coordinate_x, coordinate_y,
          pitch_result, balls_before, strikes_before, outs_before,
          velocity_mph, created_at;

-- name: GetOutsForInning :one
SELECT COALESCE(SUM(outs_recorded), 0)
FROM at_bats
WHERE game_id = $1 AND inning = $2 AND is_top_inning = $3;

-- name: GetStandings :many
SELECT t.id as team_id,
       COUNT(g.id) as games_played,
       SUM(CASE WHEN (g.home_team_id = t.id AND g.home_score > g.away_score)
                 OR (g.away_team_id = t.id AND g.away_score > g.home_score) THEN 1 ELSE 0 END) as won,
       SUM(CASE WHEN (g.home_team_id = t.id AND g.home_score < g.away_score)
                 OR (g.away_team_id = t.id AND g.away_score < g.home_score) THEN 1 ELSE 0 END) as lost
FROM teams t
LEFT JOIN games g ON (g.home_team_id = t.id OR g.away_team_id = t.id) AND g.status = 'FINAL'
GROUP BY t.id
ORDER BY won DESC, lost ASC;
```

## Modelos Generados (`internal/infrastructure/db/models.go`)

```go
// Code generated by sqlc. DO NOT EDIT.

type Game struct {
    ID            uuid.UUID
    HomeTeamID    pgtype.Text
    AwayTeamID    pgtype.Text
    Status        NullGameStatus
    HomeScore     pgtype.Int4
    AwayScore     pgtype.Int4
    CurrentInning pgtype.Int4
    IsTopInning   pgtype.Bool
    StartTime     time.Time
    CreatedAt     pgtype.Timestamptz
}

type NullGameStatus struct {
    GameStatus GameStatus
    Valid      bool
}
```

## Conversión Domain ↔ DB

```go
// internal/infrastructure/repositories/game_repository.go
func (r *GameRepository) toDomainGame(dbGame db.Game) *domain.Game {
    return &domain.Game{
        ID:            dbGame.ID,
        HomeTeamID:    dbGame.HomeTeamID.String,
        AwayTeamID:    dbGame.AwayTeamID.String,
        Status:        domain.GameStatus(dbGame.Status.GameStatus),
        HomeScore:     int(dbGame.HomeScore.Int32),
        AwayScore:     int(dbGame.AwayScore.Int32),
        CurrentInning: int(dbGame.CurrentInning.Int32),
        IsTopInning:   dbGame.IsTopInning.Bool,
        StartTime:     dbGame.StartTime,
        CreatedAt:     dbGame.CreatedAt.Time,
    }
}

func (r *GameRepository) toDBGame(domainGame *domain.Game) db.CreateGameParams {
    return db.CreateGameParams{
        HomeTeamID:    pgtype.Text{String: domainGame.HomeTeamID, Valid: true},
        AwayTeamID:    pgtype.Text{String: domainGame.AwayTeamID, Valid: true},
        Status:        db.NullGameStatus{GameStatus: db.GameStatus(domainGame.Status), Valid: true},
        HomeScore:     pgtype.Int4{Int32: int32(domainGame.HomeScore), Valid: true},
        AwayScore:     pgtype.Int4{Int32: int32(domainGame.AwayScore), Valid: true},
        CurrentInning: pgtype.Int4{Int32: int32(domainGame.CurrentInning), Valid: true},
        IsTopInning:   pgtype.Bool{Bool: domainGame.IsTopInning, Valid: true},
        StartTime:     domainGame.StartTime,
    }
}
```

## Implementación Repository

```go
// internal/infrastructure/repositories/game_repository.go
type GameRepository struct {
    queries *db.Queries
    db      *pgxpool.Pool
}

func NewGameRepository(pool *pgxpool.Pool) *GameRepository {
    return &GameRepository{
        queries: db.New(pool),
        db:      pool,
    }
}

func (r *GameRepository) GetGameByID(ctx context.Context, id uuid.UUID) (*domain.Game, error) {
    dbGame, err := r.queries.GetGameByID(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrGameNotFound
        }
        return nil, err
    }
    return r.toDomainGame(dbGame), nil
}

func (r *GameRepository) CreatePitch(ctx context.Context, pitch *domain.Pitch) (*domain.Pitch, error) {
    params := db.CreatePitchParams{
        AtBatID:       pgtype.UUID{Bytes: pitch.AtBatID, Valid: true},
        PitchNumber:   int32(pitch.PitchNumber),
        CoordinateX:   pgtype.Numeric{...}, // conversión float64 → numeric
        CoordinateY:   pgtype.Numeric{...},
        PitchResult:   pitch.PitchResult,
        BallsBefore:   int32(pitch.BallsBefore),
        StrikesBefore: int32(pitch.StrikesBefore),
        OutsBefore:    int32(pitch.OutsBefore),
        VelocityMph:   pgtype.Numeric{...}, // *float64 → numeric
    }
    dbPitch, err := r.queries.CreatePitch(ctx, params)
    return r.toDomainPitch(dbPitch), err
}
```

## Tests con Testcontainers

```go
// internal/infrastructure/repositories/game_repository_test.go
func TestGameRepository_Integration(t *testing.T) {
    // Setup: testcontainers PostgreSQL
    ctx := context.Background()
    container, _ := postgres.Run(ctx, "postgres:16-alpine")
    defer container.Terminate(ctx)
    
    connStr, _ := container.ConnectionString(ctx, "sslmode=disable")
    pool, _ := pgxpool.New(ctx, connStr)
    defer pool.Close()
    
    // Run migrations
    migrate.Migrate(connStr)
    
    repo := NewGameRepository(pool)
    
    // Test: Create → Get → Update
    game := &domain.Game{
        HomeTeamID: "ldc", AwayTeamID: "ndm",
        Status: domain.GameStatusScheduled,
        StartTime: time.Now().Add(time.Hour),
    }
    created, _ := repo.CreateGame(ctx, game)
    assert.NotNil(t, created.ID)
    
    fetched, _ := repo.GetGameByID(ctx, created.ID)
    assert.Equal(t, "ldc", fetched.HomeTeamID)
}
```

## Migraciones

```sql
-- sql/schema.sql (versión actual)
-- Ejecutar en orden:
-- 1. CREATE TYPE game_status AS ENUM (...)
-- 2. CREATE TABLE users (...)
-- 3. CREATE TABLE teams (...)
-- 4. CREATE TABLE games (...)
-- 5. CREATE TABLE at_bats (...)
-- 6. CREATE TABLE pitches (...)
-- 7. CREATE INDEXES
```

## Regenerar Código sqlc

```bash
cd backend
sqlc generate
# O con Docker:
docker run --rm -v $(pwd):/src -w /src sqlc/sqlc generate
```

## Decisiones de Diseño

| Decisión | Justificación |
|----------|---------------|
| `sqlc` + `pgx/v5` | Type-safe, performance, sin ORM overhead |
| Interfaces en `core/ports` | Dependency Inversion, testabilidad con mocks |
| Conversión explícita Domain↔DB | Control total, sin magic mapping |
| `pgtype` para nullable | Manejo correcto NULL en BD |
| Testcontainers en tests | BD real, sin mocks frágiles |

## Próximos Pasos

- [ ] Migración versionada (golang-migrate)
- [ ] Connection pooling tuning (`max_conns`, `min_conns`)
- [ ] Read replicas para queries de lectura (standings, boxscore)
- [ ] Soft deletes para auditoría