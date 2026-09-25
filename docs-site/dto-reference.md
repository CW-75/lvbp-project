# Referencia Completa de DTOs

Catálogo unificado de todos los Data Transfer Objects, Request/Response bodies, y Eventos SSE del sistema.

---

## 📥 Request DTOs (Ingreso de Datos)

### `PitchInput` - POST `/games/{gameId}/pitches`

**Ubicación**: `backend/internal/handlers/rest/ingestion_handler.go:22`

```go
type PitchInput struct {
    PitchResult string   `json:"pitchResult" validate:"required"`
    CoordinateX *float64 `json:"coordinateX,omitempty"`
    CoordinateY *float64 `json:"coordinateY,omitempty"`
    VelocityMPH *float64 `json:"velocityMPH,omitempty"`
}
```

| Campo | Tipo | Req | Descripción | Valores Válidos |
|-------|------|-----|-------------|-----------------|
| `pitchResult` | string | ✅ | Resultado del pitch según umpire | `Ball`, `Called Strike`, `Swinging Strike`, `Foul Ball`, `In Play (Out)`, `In Play (Hit)` |
| `coordinateX` | float64 | ❌ | Coordenada X zona strike | -1.5 a 1.5 (opcional) |
| `coordinateY` | float64 | ❌ | Coordenada Y zona strike | -1.5 a 1.5 (opcional) |
| `velocityMPH` | float64 | ❌ | Velocidad en mph | > 0 (opcional) |

**Ejemplo:**
```json
{
  "pitchResult": "Called Strike",
  "coordinateX": 0.15,
  "coordinateY": 0.82,
  "velocityMPH": 93.2
}
```

---

### `LoginRequest` - POST `/auth/login`

**Ubicación**: `backend/internal/auth/handlers/handler.go`

```go
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```

---

## 📤 Response DTOs (Salida de Datos)

### `PitchResponse` - Response POST `/games/{gameId}/pitches`

**Ubicación**: `backend/internal/handlers/rest/ingestion_handler.go:29`

```go
type PitchResponse struct {
    Game   *domain.Game  `json:"game"`
    AtBat  *domain.AtBat `json:"atBat"`
    Events []string      `json:"events,omitempty"`
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `game` | `Game` | Estado actualizado del juego |
| `atBat` | `AtBat` | At-bat actual (puede tener `result` si completado) |
| `events` | `string[]` | Eventos disparados: `["PitchThrown"]`, `["PitchThrown","AtBatCompleted","InningChanged"]` |

---

### `Game` - Domain Entity (usado en BoxScore, SSE)

**Ubicación**: `backend/internal/core/domain/game.go:19`

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
}
```

**GameStatus Enum:**
| Valor | Descripción |
|-------|-------------|
| `SCHEDULED` | Programado, no iniciado |
| `IN_PROGRESS` | En juego |
| `FINAL` | Finalizado |
| `POSTPONED` | Pospuesto |
| `SUSPENDED` | Suspendido |

---

### `AtBat` - Domain Entity

**Ubicación**: `backend/internal/core/domain/at_bat.go:9`

```go
type AtBat struct {
    ID           uuid.UUID
    GameID       uuid.UUID
    Inning       int
    IsTopInning  bool
    BatterID     string
    PitcherID    string
    Result       *string  // nil = en progreso
    RunsScored   int
    OutsRecorded int
    CreatedAt    time.Time
}
```

**Result Values (cuando no nil):**
`"Strikeout"`, `"Base on Balls"`, `"Hit"`, `"Out"`, `"Sacrifice"`, `"Error"`, etc.

---

### `Pitch` - Domain Entity

**Ubicación**: `backend/internal/core/domain/pitch.go:9`

```go
type Pitch struct {
    ID            uuid.UUID
    AtBatID       uuid.UUID
    PitchNumber   int
    CoordinateX   float64
    CoordinateY   float64
    PitchResult   string
    BallsBefore   int
    StrikesBefore int
    OutsBefore    int
    VelocityMPH   *float64
    CreatedAt     time.Time
}
```

---

### `TeamStandings` - GET `/standings`

**Ubicación**: `backend/internal/core/domain/standings.go:3`

```go
type TeamStandings struct {
    TeamID       int     `json:"teamId"`
    GamesPlayed  int     `json:"gamesPlayed"`
    Won          int     `json:"won"`
    Lost         int     `json:"lost"`
    Pct          float64 `json:"pct"`
    GamesBehind  float64 `json:"gamesBehind"`
}
```

---

### `AuthResponse` - POST `/auth/login`

```go
type AuthResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}

type User struct {
    ID        uuid.UUID `json:"id"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt"`
}
```

---

## 📡 Eventos SSE (Server-Sent Events)

Canal: `game:{gameId}:events` (Redis Pub/Sub)

### Formato Común
```
event: <EventType>
data: {"type":"<EventType>", ...payload}
```

### `connected` - Conexión Establecida
```json
{
  "type": "connected",
  "status": "connected",
  "gameId": "uuid"
}
```

### `PitchThrown` - Pitch Registrado
```json
{
  "type": "PitchThrown",
  "pitch": {
    "id": "uuid",
    "atBatId": "uuid",
    "pitchNumber": 3,
    "coordinateX": 0.15,
    "coordinateY": 0.82,
    "pitchResult": "Called Strike",
    "ballsBefore": 1,
    "strikesBefore": 1,
    "outsBefore": 0,
    "velocityMPH": 93.2,
    "createdAt": "2026-09-07T19:30:00Z"
  }
}
```

### `AtBatStarted` - Nuevo At-Bat
```json
{
  "type": "AtBatStarted",
  "atBat": {
    "id": "uuid",
    "gameId": "uuid",
    "inning": 3,
    "isTopInning": true,
    "batterId": "player_123",
    "pitcherId": "player_456",
    "result": null,
    "runsScored": 0,
    "outsRecorded": 0,
    "createdAt": "2026-09-07T19:30:00Z"
  }
}
```

### `AtBatCompleted` - At-Bat Terminó
```json
{
  "type": "AtBatCompleted",
  "atBat": {
    "id": "uuid",
    "gameId": "uuid",
    "inning": 3,
    "isTopInning": true,
    "batterId": "player_123",
    "pitcherId": "player_456",
    "result": "Strikeout",
    "runsScored": 0,
    "outsRecorded": 1,
    "createdAt": "2026-09-07T19:30:00Z"
  }
}
```

### `InningChanged` - Cambio de Inning (3 outs)
```json
{
  "type": "InningChanged",
  "game": { /* Game object completo */ }
}
```

### `ScoreUpdated` - Carreras Anotadas
```json
{
  "type": "ScoreUpdated",
  "game": { /* Game object completo */ },
  "runs_added": 2
}
```

---

## 🌐 Frontend Types (TypeScript)

**Ubicación**: `web/src/lib/mock-data.ts`

### `Team`
```typescript
interface Team {
  id: string;
  name: string;
  shortName: string;
  city: string;
  primaryColor: string;
  logoInitials: string;
  logoUrl?: string;
}
```

### `Game` (Frontend)
```typescript
interface Game {
  id: string;
  homeTeam: Team;
  awayTeam: Team;
  homeScore: number;
  awayScore: number;
  inning: number;
  isTopInning: boolean;
  status: "live" | "scheduled" | "final";
  scheduledAt: string;
  stadium?: string;
  homeHits?: number;
  awayHits?: number;
  homeErrors?: number;
  awayErrors?: number;
  pitcherWin?: string;
  pitcherLoss?: string;
  pitcherSave?: string;
  pitcherProbableHome?: string;
  pitcherProbableAway?: string;
}
```

### `StandingsEntry`
```typescript
interface StandingsEntry {
  team: Team;
  gamesPlayed: number;
  wins: number;
  losses: number;
  pct: string;  // ".656"
  diff: string; // "2.0" o "--"
}
```

### `NewsArticle`
```typescript
interface NewsArticle {
  id: string;
  title: string;
  excerpt: string;
  category: "Resultados" | "Fichajes" | "Equipos" | "Temporada" | "Destacado";
  imageUrl: string;
  publishedAt: string;
  isFeatured: boolean;
}
```

---

## 🔄 Mapeo Backend ↔ Frontend

| Backend (Go) | Frontend (TS) | Conversión |
|--------------|---------------|------------|
| `domain.Game` | `Game` | `status` enum → `"live"\|"scheduled"\|"final"` |
| `domain.AtBat` | (embebido en Game) | Mismo estructura |
| `domain.Pitch` | (embebido en Game) | Mismo estructura |
| `domain.TeamStandings` | `StandingsEntry` | `TeamID` → lookup `Team` mock |
| `PitchInput` | (form data) | Mismo campo a campo |

---

## 📋 Tabla Resumen de Endpoints

| Método | Endpoint | Request | Response | Auth |
|--------|----------|---------|----------|------|
| POST | `/auth/login` | `LoginRequest` | `AuthResponse` | No |
| POST | `/auth/logout` | - | `204` | Sí |
| GET | `/auth/me` | - | `User` | Sí |
| POST | `/games/{id}/pitches` | `PitchInput` | `PitchResponse` | Sí |
| GET | `/games/{id}/boxscore` | - | `Game` | No |
| GET | `/standings` | - | `TeamStandings[]` | No |
| GET | `/games/{id}/stream` | - | SSE Events | No* |

*SSE auth opcional para juegos públicos

---

## 🎯 Validaciones Comunes

```go
// PitchResult enum válido
var ValidPitchResults = []string{
    "Ball",
    "Called Strike",
    "Swinging Strike",
    "Foul Ball",
    "In Play (Out)",
    "In Play (Hit)",
}

// Coordenadas strike zone
const (
    StrikeZoneXMin = -1.5
    StrikeZoneXMax = 1.5
    StrikeZoneYMin = -1.5
    StrikeZoneYMax = 1.5
)
```

---

## 📝 Changelog DTOs

| Versión | Fecha | Cambios |
|---------|-------|---------|
| 1.0 | 2026-09-24 | Versión inicial - todas las fases documentadas |