# Fase 7: Handler Boxscore (Box Score)

## Resumen

Endpoint REST para obtener el detalle completo de un juego: score por inning, at-bats, pitches, líneas de bateo/picheo.

## Endpoint

```
GET /games/{gameId}/boxscore
Accept: application/json
```

## Handler (`internal/handlers/rest/boxscore_handler.go`)

```go
type BoxscoreHandler struct {
    repo ports.GameRepository
}

func (h *BoxscoreHandler) GetBoxscore(w http.ResponseWriter, r *http.Request) {
    gameIDStr := chi.URLParam(r, "gameId")
    gameID, err := uuid.Parse(gameIDStr)
    if err != nil {
        http.Error(w, "Invalid game ID format", http.StatusBadRequest)
        return
    }

    game, err := h.repo.GetGameByID(r.Context(), gameID)
    if err != nil {
        http.Error(w, "Game not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(game)
}
```

## Repository Interface

```go
type GameRepository interface {
    GetGameByID(ctx context.Context, id uuid.UUID) (*domain.Game, error)
    GetAtBatsByGame(ctx context.Context, gameID uuid.UUID) ([]*domain.AtBat, error)
    GetPitchesByGame(ctx context.Context, gameID uuid.UUID) ([]*domain.Pitch, error)
}
```

## Response: `domain.Game` (con relaciones cargadas)

```go
type Game struct {
    ID            uuid.UUID
    HomeTeamID    string
    AwayTeamID    string
    Status        GameStatus
    HomeScore     int
    AwayScore     int
    CurrentInning int
    IsTopInning   bool
    StartTime     time.Time
    CreatedAt     time.Time
    
    // Cargados por GetGameByID (eager loading)
    AtBats  []*AtBat
    Pitches []*Pitch
}
```

### Ejemplo Response

```json
{
  "id": "uuid",
  "homeTeamId": "ldc",
  "awayTeamId": "ndm",
  "status": "FINAL",
  "homeScore": 5,
  "awayScore": 3,
  "currentInning": 9,
  "isTopInning": false,
  "startTime": "2026-09-07T19:00:00-04:00",
  "createdAt": "2026-09-01T10:00:00Z",
  "atBats": [
    {
      "id": "uuid",
      "gameId": "uuid",
      "inning": 1,
      "isTopInning": true,
      "batterId": "player_001",
      "pitcherId": "player_020",
      "result": "Strikeout",
      "runsScored": 0,
      "outsRecorded": 1,
      "createdAt": "2026-09-07T19:05:00Z"
    }
  ],
  "pitches": [
    {
      "id": "uuid",
      "atBatId": "uuid",
      "pitchNumber": 1,
      "coordinateX": -0.3,
      "coordinateY": 0.8,
      "pitchResult": "Ball",
      "ballsBefore": 0,
      "strikesBefore": 0,
      "outsBefore": 0,
      "velocityMPH": 91.2,
      "createdAt": "2026-09-07T19:05:00Z"
    }
  ]
}
```

## Frontend Integration (TanStack Query)

```typescript
// web/src/hooks/useBoxscore.ts
export function useBoxscore(gameId: string) {
  return useQuery({
    queryKey: ['boxscore', gameId],
    queryFn: async (): Promise<Game> => {
      const res = await fetch(`/api/games/${gameId}/boxscore`)
      if (!res.ok) throw new Error('Failed to fetch boxscore')
      return res.json()
    },
    enabled: !!gameId,
    staleTime: 5 * 60 * 1000, // 5 min para juegos finalizados
  })
}
```

## Estructura de Datos para UI

### Score by Inning (Line Score)

```typescript
// Derivado de Game + AtBats
interface LineScore {
  inning: number
  awayRuns: number
  homeRuns: number
}

// Cálculo:
// 1. Agrupar AtBats por inning + isTopInning
// 2. Sumar runsScored por grupo
// 3. Rellenar innings vacíos con 0
```

### Batting Lines

```typescript
interface BattingLine {
  playerId: string
  playerName: string
  teamId: string
  ab: number    // At-bats oficiales
  r: number     // Runs
  h: number     // Hits
  rbi: number   // RBI
  bb: number    // Walks
  k: number     // Strikeouts
  avg: number   // BA
}
```

### Pitching Lines

```typescript
interface PitchingLine {
  playerId: string
  playerName: string
  teamId: string
  ip: number    // Innings pitched (outs/3)
  h: number     // Hits allowed
  r: number     // Runs
  er: number    // Earned runs
  bb: number    // Walks
  k: number     // Strikeouts
  hr: number    // Home runs
  era: number   // ERA
}
```

## Queries SQL (sqlc)

```sql
-- queries.sql
-- name: GetGameWithDetails :one
SELECT g.*, 
       json_agg(DISTINCT ab.*) FILTER (WHERE ab.id IS NOT NULL) as at_bats,
       json_agg(DISTINCT p.*) FILTER (WHERE p.id IS NOT NULL) as pitches
FROM games g
LEFT JOIN at_bats ab ON ab.game_id = g.id
LEFT JOIN pitches p ON p.at_bat_id = ab.id
WHERE g.id = $1
GROUP BY g.id;
```

## Performance

| Métrica | Estrategia |
|---------|------------|
| Latencia | Query única con JOINs + JSON aggregation |
| Cache | TanStack Query 5min (juegos FINAL), 30s (IN_PROGRESS) |
| Payload | ~50KB típico (game + 50 at-bats + 150 pitches) |

## Tests

```go
// internal/handlers/rest/boxscore_handler_test.go
func TestGetBoxscore_ReturnsGameWithAtBatsAndPitches(t *testing.T) {
    // Setup: game con 3 at-bats, 10 pitches
    // Act: GET /games/{id}/boxscore
    // Assert: 200, Game con AtBats y Pitches poblados
}

func TestGetBoxscore_NotFound_Returns404(t *testing.T) {
    // Act: GET /games/invalid-id/boxscore
    // Assert: 404 "Game not found"
}
```

## Próximos Pasos

- [ ] Endpoint `/games/{id}/linescore` (solo line score, payload pequeño)
- [ ] Endpoint `/games/{id}/playbyplay` (solo at-bats ordenados)
- [ ] Box score PDF export
- [ ] Statcast-style metrics (exit velo, launch angle - requiere datos tracking)