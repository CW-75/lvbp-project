---
layout: page
title: API Reference
---

# API Reference (OpenAPI 3.0)

Documentación interactiva de la API REST generada desde código Go con `swaggo/swag`.

## 🔗 Acceso Rápido

- **Swagger UI**: [/api/swagger](/api/swagger) (cuando servidor docs corriendo)
- **Scalar API Reference**: [/api/scalar](/api/scalar) (moderno, tema oscuro)
- **OpenAPI Spec**: [/api/openapi.yaml](/api/openapi.yaml) (raw YAML)

## 📋 Endpoints Resumidos

### Autenticación
| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/auth/login` | Login scorekeeper |
| POST | `/auth/logout` | Logout (revoca sesión) |
| GET | `/auth/me` | Usuario actual |

### Juegos - Ingesta (Scorekeeper)
| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/games/{gameId}/pitches` | Registrar pitch |
| POST | `/games/{gameId}/at-bats` | Iniciar at-bat (TODO) |

### Juegos - Consulta (Público)
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/games/{gameId}/boxscore` | Box score completo |
| GET | `/games/{gameId}/linescore` | Line score solo (TODO) |
| GET | `/games/{gameId}/playbyplay` | Play-by-play (TODO) |
| GET | `/games/{gameId}/stream` | **SSE** Tiempo real |

### Standings
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/standings` | Tabla de posiciones |

## 🛠 Generación Local

```bash
# Backend: generar openapi.yaml
cd backend
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o ../docs-site/public/api --parseDependency --parseInternal

# Frontend: generar typedoc
cd web
npx typedoc --out ../docs-site/public/typedoc --entryPoints src/lib
```

## 📦 Estructura OpenAPI

```yaml
openapi: 3.0.3
info:
  title: LVBP GameCast API
  version: 1.0.0
  description: Live Baseball GameCast & Standings Engine
servers:
  - url: http://localhost:8080
    description: Development
  - url: https://api.lvbp.example.com
    description: Production

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
  schemas:
    PitchInput:
      type: object
      required: [pitchResult]
      properties:
        pitchResult:
          type: string
          enum: [Ball, "Called Strike", "Swinging Strike", "Foul Ball", "In Play (Out)", "In Play (Hit)"]
        coordinateX:
          type: number
          format: float
        coordinateY:
          type: number
          format: float
        velocityMPH:
          type: number
          format: float
    PitchResponse:
      type: object
      properties:
        game:
          $ref: '#/components/schemas/Game'
        atBat:
          $ref: '#/components/schemas/AtBat'
        events:
          type: array
          items:
            type: string
    Game:
      type: object
      properties:
        id:
          type: string
          format: uuid
        homeTeamId:
          type: string
        awayTeamId:
          type: string
        status:
          type: string
          enum: [SCHEDULED, IN_PROGRESS, FINAL, POSTPONED, SUSPENDED]
        homeScore:
          type: integer
        awayScore:
          type: integer
        currentInning:
          type: integer
        isTopInning:
          type: boolean
        startTime:
          type: string
          format: date-time
        createdAt:
          type: string
          format: date-time
    AtBat:
      type: object
      properties:
        id:
          type: string
          format: uuid
        gameId:
          type: string
          format: uuid
        inning:
          type: integer
        isTopInning:
          type: boolean
        batterId:
          type: string
        pitcherId:
          type: string
        result:
          type: string
          nullable: true
        runsScored:
          type: integer
        outsRecorded:
          type: integer
        createdAt:
          type: string
          format: date-time
    TeamStandings:
      type: object
      properties:
        teamId:
          type: integer
        gamesPlayed:
          type: integer
        won:
          type: integer
        lost:
          type: integer
        pct:
          type: number
          format: float
        gamesBehind:
          type: number
          format: float
```

## 🔐 Autenticación

### JWT Bearer Token
```
Authorization: Bearer <token>
```

### Cookie HttpOnly (Alternativa)
```
Cookie: auth_token=<token>
```

**Obtener token:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"scorekeeper@lvbp.com","password":"password123"}'
```

## 📡 Server-Sent Events (SSE)

### Conexión
```bash
curl -N -H "Accept: text/event-stream" \
  http://localhost:8080/games/{gameId}/stream
```

### Formato Eventos
```
event: PitchThrown
data: {"type":"PitchThrown","pitch":{...}}

event: AtBatCompleted
data: {"type":"AtBatCompleted","atBat":{...}}

event: InningChanged
data: {"type":"InningChanged","game":{...}}

event: ScoreUpdated
data: {"type":"ScoreUpdated","game":{...},"runs_added":2}
```

### JavaScript (EventSource)
```javascript
const es = new EventSource(`/games/${gameId}/stream`)

es.addEventListener('PitchThrown', (e) => {
  const event = JSON.parse(e.data)
  updatePitchDisplay(event.pitch)
})

es.addEventListener('ScoreUpdated', (e) => {
  const event = JSON.parse(e.data)
  updateScoreboard(event.game)
})
```

## ⚠️ Códigos de Error

| Código | Significado | Cuándo |
|--------|-------------|--------|
| 200 | OK | Éxito |
| 201 | Created | Recurso creado |
| 204 | No Content | Éxito sin body (logout) |
| 400 | Bad Request | JSON inválido, validación fallida |
| 401 | Unauthorized | Token inválido/expirado |
| 404 | Not Found | Recurso no existe |
| 409 | Conflict | Estado inválido (ej. juego no IN_PROGRESS) |
| 422 | Unprocessable | Validación de negocio fallida |
| 500 | Internal Server Error | Error inesperado |

## 📊 Rate Limiting (Planeado)

| Endpoint | Límite | Ventana |
|----------|--------|---------|
| `/auth/login` | 5 req | 1 min |
| `/games/{id}/pitches` | 30 req | 1 min |
| `/standings` | 60 req | 1 min |
| `/games/{id}/boxscore` | 30 req | 1 min |

## 🧪 Ejemplos cURL

### Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"scorekeeper@lvbp.com","password":"securePass123"}'
```

### Registrar Pitch
```bash
curl -X POST http://localhost:8080/games/abc-123/pitches \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"pitchResult":"Swinging Strike","coordinateX":0.2,"coordinateY":0.7,"velocityMPH":95.1}'
```

### Obtener Box Score
```bash
curl http://localhost:8080/games/abc-123/boxscore
```

### Obtener Standings
```bash
curl http://localhost:8080/standings
```

### Stream SSE
```bash
curl -N http://localhost:8080/games/abc-123/stream \
  -H "Accept: text/event-stream"
```