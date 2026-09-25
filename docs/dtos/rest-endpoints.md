# Endpoints REST - DTOs

---

## 1. GET /standings

Obtiene la tabla de posiciones de la liga.

### Response: `TeamStandings[]` (200 OK)

```json
[
  {
    "teamId": 1,
    "gamesPlayed": 15,
    "won": 10,
    "lost": 5,
    "pct": 0.667,
    "gamesBehind": 0.0
  },
  {
    "teamId": 2,
    "gamesPlayed": 14,
    "won": 8,
    "lost": 6,
    "pct": 0.571,
    "gamesBehind": 1.5
  }
]
```

### Campos `TeamStandings`

| Campo | Tipo | Obligatorio | Descripción |
|-------|------|-------------|-------------|
| `teamId` | `integer` | Sí | Identificador único del equipo |
| `gamesPlayed` | `integer` | Sí | Total de juegos disputados |
| `won` | `integer` | Sí | Juegos ganados |
| `lost` | `integer` | Sí | Juegos perdidos |
| `pct` | `number` (float) | Sí | Porcentaje de victorias (won / gamesPlayed), 3 decimales |
| `gamesBehind` | `number` (float) | Sí | Juegos de diferencia vs líder, 1 decimal |

---

## 2. GET /games/{gameId}/boxscore

Obtiene el boxscore completo de un juego.

### Parámetros de Ruta

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `gameId` | `string` (UUID) | Identificador único del juego |

### Response: `Game` (200 OK)

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "homeTeamId": "LARA",
  "awayTeamId": "MAGA",
  "status": "IN_PROGRESS",
  "homeScore": 3,
  "awayScore": 2,
  "currentInning": 5,
  "isTopInning": true,
  "startTime": "2025-10-15T23:00:00.000Z",
  "createdAt": "2025-10-15T22:45:00.000Z"
}
```

### Response Error (404 Not Found)

```json
{
  "error": "Game not found"
}
```

### Campos `Game`

| Campo | Tipo | Obligatorio | Descripción |
|-------|------|-------------|-------------|
| `id` | `string` (UUID) | Sí | Identificador único del juego |
| `homeTeamId` | `string` | Sí | Código del equipo local (ej: "LARA", "MAGA", "CARI", "LEON") |
| `awayTeamId` | `string` | Sí | Código del equipo visitante |
| `status` | `string` (enum) | Sí | `SCHEDULED`, `IN_PROGRESS`, `FINAL`, `POSTPONED`, `SUSPENDED` |
| `homeScore` | `integer` | Sí | Carreras del equipo local |
| `awayScore` | `integer` | Sí | Carreras del equipo visitante |
| `currentInning` | `integer` | Sí | Inning actual (1-9+ para extra innings) |
| `isTopInning` | `boolean` | Sí | `true` = parte alta (batea visitante), `false` = parte baja (batea local) |
| `startTime` | `string` (ISO 8601) | Sí | Hora programada de inicio |
| `createdAt` | `string` (ISO 8601) | Sí | Timestamp de creación del registro |

---

## 3. POST /games/{gameId}/pitches

Registra un lanzamiento y procesa la lógica del motor de béisbol.

### Parámetros de Ruta

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `gameId` | `string` (UUID) | Identificador único del juego |

### Request: `PitchInput`

```json
{
  "pitchResult": "Swinging Strike",
  "coordinateX": 0.75,
  "coordinateY": -0.25,
  "velocityMPH": 98.5
}
```

### Request Mínimo

```json
{
  "pitchResult": "Ball"
}
```

### Campos `PitchInput`

| Campo | Tipo | Obligatorio | Descripción | Validaciones |
|-------|------|-------------|-------------|--------------|
| `pitchResult` | `string` | **Sí** | Resultado del lanzamiento | Ver valores válidos abajo |
| `coordinateX` | `number` (float) | No | Coordenada X zona strike (-1.0 a 1.0) | Rango [-1.0, 1.0] |
| `coordinateY` | `number` (float) | No | Coordenada Y zona strike (-1.0 a 1.0) | Rango [-1.0, 1.0] |
| `velocityMPH` | `number` (float) | No | Velocidad del lanzamiento en mph | > 0, típicamente 70-105 |

### Valores Válidos para `pitchResult`

| Valor | Descripción | Efecto en Conteo |
|-------|-------------|------------------|
| `"Ball"` | Bola mala | Ball +1 (si 4 → Base on Balls) |
| `"Called Strike"` | Strike cantado | Strike +1 (si 3 → Strikeout) |
| `"Swinging Strike"` | Strike tirando | Strike +1 (si 3 → Strikeout) |
| `"Foul Ball"` | Foul | Strike +1 solo si strikes < 2 |
| `"In Play (Out)"` | En juego - Out | At-bat completado, out +1 |
| `"In Play (Hit)"` | En juego - Hit | At-bat completado, posible carrera |

### Response: `PitchResponse` (200 OK) - At-bat completado

```json
{
  "game": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "homeTeamId": "LARA",
    "awayTeamId": "MAGA",
    "status": "IN_PROGRESS",
    "homeScore": 3,
    "awayScore": 2,
    "currentInning": 5,
    "isTopInning": true,
    "startTime": "2025-10-15T23:00:00.000Z",
    "createdAt": "2025-10-15T22:45:00.000Z"
  },
  "atBat": {
    "id": "b2c3d4e5-f6a7-8901-bcde-f23456789012",
    "gameId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "inning": 5,
    "isTopInning": true,
    "batterId": "PLAYER_001",
    "pitcherId": "PLAYER_045",
    "result": "Strikeout",
    "runsScored": 0,
    "outsRecorded": 1,
    "createdAt": "2025-10-15T23:15:30.000Z"
  },
  "events": ["PitchThrown", "AtBatCompleted", "InningChanged"]
}
```

### Response: `PitchResponse` - At-bat continúa

```json
{
  "game": { ... },
  "atBat": {
    "id": "b2c3d4e5-f6a7-8901-bcde-f23456789012",
    "gameId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "inning": 5,
    "isTopInning": true,
    "batterId": "PLAYER_001",
    "pitcherId": "PLAYER_045",
    "result": null,
    "runsScored": 0,
    "outsRecorded": 0,
    "createdAt": "2025-10-15T23:15:30.000Z"
  },
  "events": ["PitchThrown"]
}
```

### Response Errores

| Código | Causa | Body |
|--------|-------|------|
| 400 | `gameId` inválido | `{"error": "Invalid game ID format"}` |
| 400 | JSON inválido | `{"error": "Invalid JSON payload"}` |
| 400 | `pitchResult` no válido | `{"error": "Invalid pitchResult value"}` |
| 409 | Juego no `IN_PROGRESS` | `{"error": "game is not in progress"}` |
| 409 | At-bat ya completado | `{"error": "current at-bat is already completed"}` |
| 500 | Error interno | `{"error": "Internal server error"}` |

### Campos `PitchResponse`

| Campo | Tipo | Obligatorio | Descripción |
|-------|------|-------------|-------------|
| `game` | `Game` | Sí | Estado actualizado del juego |
| `atBat` | `AtBat` | Sí | Estado actualizado del turno al bate |
| `events` | `string[]` | Sí | Eventos generados (ver `sse-events.md`) |

### Campos `AtBat`

| Campo | Tipo | Obligatorio | Descripción |
|-------|------|-------------|-------------|
| `id` | `string` (UUID) | Sí | Identificador único del at-bat |
| `gameId` | `string` (UUID) | Sí | Referencia al juego |
| `inning` | `integer` | Sí | Inning en que ocurrió |
| `isTopInning` | `boolean` | Sí | Parte alta/baja del inning |
| `batterId` | `string` | Sí | ID del bateador |
| `pitcherId` | `string` | Sí | ID del lanzador |
| `result` | `string` \| `null` | Sí | `null` si en curso, o `"Strikeout"`, `"Base on Balls"`, `"Out"`, `"Hit"` |
| `runsScored` | `integer` | Sí | Carreras anotadas en este at-bat |
| `outsRecorded` | `integer` | Sí | Outs registrados (0 o 1) |
| `createdAt` | `string` (ISO 8601) | Sí | Timestamp de creación |

---

## Referencias

- [Eventos SSE →](sse-events.md)
- [Reglas de negocio →](business-rules.md)
- [Ejemplos de flujo →](flow-examples.md)
- [Tipos TypeScript →](typescript-types.md)