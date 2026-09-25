# Eventos SSE (Server-Sent Events) - DTOs

---

## Endpoint

```
GET /games/{gameId}/stream
```

## Headers de Respuesta

```
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
X-Accel-Buffering: no
```

## Formato de Mensaje

Cada mensaje sigue el estándar SSE:

```
event: <tipo_evento>
data: <json_payload>

```

> **Nota:** Doble salto de línea (`\n\n`) al final de cada mensaje.

---

## 1. `connected`

Enviado inmediatamente tras establecer la conexión.

### Ejemplo

```
event: connected
data: {"status":"connected","gameId":"a1b2c3d4-e5f6-7890-abcd-ef1234567890"}
```

### Payload

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `status` | `string` | Siempre `"connected"` |
| `gameId` | `string` (UUID) | ID del juego suscrito |

---

## 2. `PitchThrown`

Emitido cada vez que se registra un lanzamiento válido.

### Ejemplo

```
event: PitchThrown
data: {"type":"PitchThrown","pitch":{"id":"c3d4e5f6-a7b8-9012-cdef-345678901234","atBatId":"b2c3d4e5-f6a7-8901-bcde-f23456789012","pitchNumber":3,"coordinateX":0.75,"coordinateY":-0.25,"pitchResult":"Swinging Strike","ballsBefore":2,"strikesBefore":1,"outsBefore":1,"velocityMPH":98.5,"createdAt":"2025-10-15T23:16:45.000Z"}}
```

### Payload

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `type` | `string` | Siempre `"PitchThrown"` |
| `pitch` | `object` | Objeto `Pitch` completo |

### Campos `Pitch`

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | `string` (UUID) | ID único del lanzamiento |
| `atBatId` | `string` (UUID) | ID del at-bat al que pertenece |
| `pitchNumber` | `integer` | Número de lanzamiento en el at-bat (1, 2, 3...) |
| `coordinateX` | `number` (float) | Coordenada X zona strike |
| `coordinateY` | `number` (float) | Coordenada Y zona strike |
| `pitchResult` | `string` | Resultado (ver `rest-endpoints.md`) |
| `ballsBefore` | `integer` | Bolas antes de este lanzamiento |
| `strikesBefore` | `integer` | Strikes antes de este lanzamiento |
| `outsBefore` | `integer` | Outs en el inning antes de este lanzamiento |
| `velocityMPH` | `number` (float) \| `null` | Velocidad en mph |
| `createdAt` | `string` (ISO 8601) | Timestamp del lanzamiento |

---

## 3. `AtBatStarted`

Emitido cuando inicia un nuevo turno al bate.

### Ejemplo

```
event: AtBatStarted
data: {"type":"AtBatStarted","atBat":{"id":"b2c3d4e5-f6a7-8901-bcde-f23456789012","gameId":"a1b2c3d4-e5f6-7890-abcd-ef1234567890","inning":5,"isTopInning":true,"batterId":"PLAYER_001","pitcherId":"PLAYER_045","result":null,"runsScored":0,"outsRecorded":0,"createdAt":"2025-10-15T23:15:30.000Z"}}
```

### Payload

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `type` | `string` | Siempre `"AtBatStarted"` |
| `atBat` | `object` | Objeto `AtBat` completo (ver `rest-endpoints.md`) |

---

## 4. `AtBatCompleted`

Emitido cuando finaliza un turno al bate.

### Ejemplo

```
event: AtBatCompleted
data: {"type":"AtBatCompleted","atBat":{"id":"b2c3d4e5-f6a7-8901-bcde-f23456789012","gameId":"a1b2c3d4-e5f6-7890-abcd-ef1234567890","inning":5,"isTopInning":true,"batterId":"PLAYER_001","pitcherId":"PLAYER_045","result":"Strikeout","runsScored":0,"outsRecorded":1,"createdAt":"2025-10-15T23:15:30.000Z"}}
```

### Payload

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `type` | `string` | Siempre `"AtBatCompleted"` |
| `atBat` | `object` | Objeto `AtBat` con `result` poblado |

### Valores de `atBat.result`

| Valor | Significado |
|-------|-------------|
| `"Strikeout"` | Ponche (3 strikes) |
| `"Base on Balls"` | Base por bolas (4 bolas) |
| `"Out"` | Out en juego |
| `"Hit"` | Hit en juego |

---

## 5. `InningChanged`

Emitido cuando cambian de mitad de inning o de inning (3 outs).

### Ejemplo

```
event: InningChanged
data: {"type":"InningChanged","game":{"id":"a1b2c3d4-e5f6-7890-abcd-ef1234567890","homeTeamId":"LARA","awayTeamId":"MAGA","status":"IN_PROGRESS","homeScore":3,"awayScore":2,"currentInning":5,"isTopInning":false,"startTime":"2025-10-15T23:00:00.000Z","createdAt":"2025-10-15T22:45:00.000Z"}}
```

### Payload

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `type` | `string` | Siempre `"InningChanged"` |
| `game` | `object` | Objeto `Game` actualizado (ver `rest-endpoints.md`) |

### Transiciones Típicas

| Antes | Después | Descripción |
|-------|---------|-------------|
| `inning: 5, isTopInning: true` | `inning: 5, isTopInning: false` | Cambio a parte baja |
| `inning: 5, isTopInning: false` | `inning: 6, isTopInning: true` | Nuevo inning, parte alta |

---

## 6. `ScoreUpdated`

Emitido cuando se anotan carreras.

### Ejemplo

```
event: ScoreUpdated
data: {"type":"ScoreUpdated","game":{"id":"a1b2c3d4-e5f6-7890-abcd-ef1234567890","homeTeamId":"LARA","awayTeamId":"MAGA","status":"IN_PROGRESS","homeScore":3,"awayScore":4,"currentInning":5,"isTopInning":true,"startTime":"2025-10-15T23:00:00.000Z","createdAt":"2025-10-15T22:45:00.000Z"},"runs_added":2}
```

### Payload

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `type` | `string` | Siempre `"ScoreUpdated"` |
| `game` | `object` | Objeto `Game` con marcador actualizado |
| `runs_added` | `integer` | Carreras anotadas en esta jugada (≥1) |

### Lógica de Anotación

- **Top inning (`isTopInning: true`)** → Suman a `awayScore` (batea visitante)
- **Bottom inning (`isTopInning: false`)** → Suman a `homeScore` (batea local)

---

## Resumen de Tipos de Evento

| Evento | `type` en payload | Entidad Principal | Cuándo se emite |
|--------|-------------------|-------------------|-----------------|
| `connected` | — | `{status, gameId}` | Conexión SSE establecida |
| `PitchThrown` | `"PitchThrown"` | `Pitch` | Cada lanzamiento registrado |
| `AtBatStarted` | `"AtBatStarted"` | `AtBat` | Inicio de turno al bate |
| `AtBatCompleted` | `"AtBatCompleted"` | `AtBat` | Fin de turno (resultado definido) |
| `InningChanged` | `"InningChanged"` | `Game` | 3 outs → cambio mitad/inning |
| `ScoreUpdated` | `"ScoreUpdated"` | `Game` + `runs_added` | Carreras anotadas |

---

## Referencias

- [Endpoints REST →](rest-endpoints.md)
- [Reglas de negocio →](business-rules.md)
- [Ejemplos de flujo →](flow-examples.md)
- [Tipos TypeScript →](typescript-types.md)