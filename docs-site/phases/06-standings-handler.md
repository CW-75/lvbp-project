# Fase 6: Handler Standings (Clasificación)

## Resumen

Endpoint REST para obtener la tabla de posiciones actualizada. Consulta agregada calculada en BD.

## Endpoint

```
GET /standings
Accept: application/json
```

## Handler (`internal/handlers/rest/standings_handler.go`)

```go
type StandingsHandler struct {
    repo ports.StandingsRepository
}

func (h *StandingsHandler) GetStandings(w http.ResponseWriter, r *http.Request) {
    standings, err := h.repo.GetStandings(r.Context())
    if err != nil {
        http.Error(w, "Failed to retrieve standings", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(standings)
}
```

## Repository Interface (`internal/core/ports/repositories.go`)

```go
type StandingsRepository interface {
    GetStandings(ctx context.Context) ([]*domain.TeamStandings, error)
}
```

## Implementación PostgreSQL (`internal/infrastructure/repositories/game_repository.go`)

```go
func (r *GameRepository) GetStandings(ctx context.Context) ([]*domain.TeamStandings, error) {
    // Query calculada en BD usando vista o CTE
    rows, err := r.db.Query(ctx, `
        SELECT 
            t.id as team_id,
            COUNT(g.id) as games_played,
            SUM(CASE WHEN (g.home_team_id = t.id AND g.home_score > g.away_score) 
                      OR (g.away_team_id = t.id AND g.away_score > g.home_score) THEN 1 ELSE 0 END) as won,
            SUM(CASE WHEN (g.home_team_id = t.id AND g.home_score < g.away_score) 
                      OR (g.away_team_id = t.id AND g.away_score < g.home_score) THEN 1 ELSE 0 END) as lost
        FROM teams t
        LEFT JOIN games g ON (g.home_team_id = t.id OR g.away_team_id = t.id) AND g.status = 'FINAL'
        GROUP BY t.id
        ORDER BY won DESC, lost ASC
    `)
    // ... scan a []domain.TeamStandings
}
```

## DTO Response: `TeamStandings` (`core/domain/standings.go`)

```go
type TeamStandings struct {
    TeamID       int     `json:"teamId"`
    GamesPlayed  int     `json:"gamesPlayed"`
    Won          int     `json:"won"`
    Lost         int     `json:"lost"`
    Pct          float64 `json:"pct"`           // Won / GamesPlayed
    GamesBehind  float64 `json:"gamesBehind"`   // Diferencia vs líder
}
```

### Ejemplo Response

```json
[
  {
    "teamId": 1,
    "gamesPlayed": 32,
    "won": 21,
    "lost": 11,
    "pct": 0.656,
    "gamesBehind": 0
  },
  {
    "teamId": 2,
    "gamesPlayed": 32,
    "won": 19,
    "lost": 13,
    "pct": 0.594,
    "gamesBehind": 2.0
  }
]
```

## Frontend Integration (TanStack Query)

```typescript
// web/src/hooks/useStandings.ts
import { useQuery } from '@tanstack/react-query'

interface TeamStandings {
  teamId: number
  gamesPlayed: number
  won: number
  lost: number
  pct: number
  gamesBehind: number
}

export function useStandings() {
  return useQuery({
    queryKey: ['standings'],
    queryFn: async (): Promise<TeamStandings[]> => {
      const res = await fetch('/api/standings')
      if (!res.ok) throw new Error('Failed to fetch standings')
      return res.json()
    },
    staleTime: 30_000, // 30s cache
    refetchOnWindowFocus: false,
  })
}
```

## Cálculo de Games Behind (GB)

```
GB = ((Leader_Wins - Team_Wins) + (Team_Losses - Leader_Losses)) / 2
```

Ejemplo:
- Líder: 21-11 (.656)
- Equipo: 19-13 (.594)
- GB = ((21-19) + (13-11)) / 2 = (2 + 2) / 2 = 2.0

## Performance

| Métrica | Objetivo | Implementación |
|---------|----------|----------------|
| Latencia p95 | < 100ms | Índice en `games(status)`, vista materializada |
| Cache | 30s | TanStack Query `staleTime: 30_000` |
| Frecuencia actualización | Tiempo real | SSE `ScoreUpdated` → `invalidateQueries(['standings'])` |

## Invalidación de Cache (Tiempo Real)

```typescript
// En componente que usa SSE
const queryClient = useQueryClient()

useEffect(() => {
  const unsubscribe = useGameStore.subscribe(
    (state) => state.lastScoreUpdate,
    (lastScoreUpdate) => {
      if (lastScoreUpdate) {
        queryClient.invalidateQueries({ queryKey: ['standings'] })
      }
    }
  )
  return unsubscribe
}, [queryClient])
```

## Tests

```go
// internal/handlers/rest/standings_handler_test.go
func TestGetStandings_ReturnsOrderedStandings(t *testing.T) {
    // Setup: 3 equipos con diferentes records
    // Act: GET /standings
    // Assert: 200, array ordenado por PCT desc, GB calculado
}

func TestGetStandings_EmptySeason_ReturnsZeros(t *testing.T) {
    // Setup: BD vacía
    // Assert: 200, [] (o equipos con 0-0)
}
```

## Próximos Pasos

- [ ] Vista materializada `standings_view` refrescada por trigger en `games`
- [ ] Filtros: `/standings?division=central&season=2026`
- [ ] Estadísticas avanzadas: racha, home/away, vs division
- [ ] Paginación para ligas grandes