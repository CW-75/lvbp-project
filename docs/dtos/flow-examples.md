# Ejemplos de Flujos Completos

Esta sección documenta secuencias típicas de interacción entre el cliente (Scorekeeper/Frontend) y la API, mostrando requests, responses REST y eventos SSE emitidos en orden.

---

## 1. Strikeout en 3 Lanzamientos

### Contexto
- Juego en progreso, inning 5, parte alta (batea visitante)
- Cuenta inicial: 0 bolas, 0 strikes, 1 out
- Bateador: `PLAYER_001`, Lanzador: `PLAYER_045`

### Lanzamiento 1: Called Strike

**Request:**
```http
POST /games/a1b2c3d4-e5f6-7890-abcd-ef1234567890/pitches
Content-Type: application/json

{"pitchResult": "Called Strike"}
```

**Response (200 OK):**
```json
{
  "game": { "id": "...", "status": "IN_PROGRESS", "homeScore": 3, "awayScore": 2, "currentInning": 5, "isTopInning": true },
  "atBat": { "id": "...", "result": null, "runsScored": 0, "outsRecorded": 0 },
  "events": ["PitchThrown"]
}
```

**SSE Emitido:**
```
event: PitchThrown
data: {"type":"PitchThrown","pitch":{"pitchNumber":1,"pitchResult":"Called Strike","ballsBefore":0,"strikesBefore":0,"outsBefore":1,...}}
```

---

### Lanzamiento 2: Swinging Strike

**Request:**
```http
POST /games/a1b2c3d4-e5f6-7890-abcd-ef1234567890/pitches
Content-Type: application/json

{"pitchResult": "Swinging Strike"}
```

**Response (200 OK):**
```json
{
  "game": { ... },
  "atBat": { "result": null, "runsScored": 0, "outsRecorded": 0 },
  "events": ["PitchThrown"]
}
```

**SSE Emitido:**
```
event: PitchThrown
data: {"type":"PitchThrown","pitch":{"pitchNumber":2,"pitchResult":"Swinging Strike","ballsBefore":0,"strikesBefore":1,...}}
```

---

### Lanzamiento 3: Called Strike (Strikeout)

**Request:**
```http
POST /games/a1b2c3d4-e5f6-7890-abcd-ef1234567890/pitches
Content-Type: application/json

{"pitchResult": "Called Strike"}
```

**Response (200 OK):**
```json
{
  "game": { "id": "...", "status": "IN_PROGRESS", "homeScore": 3, "awayScore": 2, "currentInning": 5, "isTopInning": true },
  "atBat": { "id": "...", "result": "Strikeout", "runsScored": 0, "outsRecorded": 1 },
  "events": ["PitchThrown", "AtBatCompleted", "InningChanged"]
}
```

**SSE Emitidos (en orden):**
```
event: PitchThrown
data: {"type":"PitchThrown","pitch":{"pitchNumber":3,"pitchResult":"Called Strike","ballsBefore":0,"strikesBefore":2,...}}

event: AtBatCompleted
data: {"type":"AtBatCompleted","atBat":{"result":"Strikeout","outsRecorded":1,...}}

event: InningChanged
data: {"type":"InningChanged","game":{"currentInning":5,"isTopInning":false,"homeScore":3,"awayScore":2,...}}
```

> **Nota:** `InningChanged` solo se emite si era el 3er out del inning. En este ejemplo era el 2do out, así que NO se emitiría `InningChanged`.

---

## 2. Base on Balls (Walk) en 4 Lanzamientos

### Contexto
- Mismo juego, nuevo at-bat, cuenta 0-0

### Lanzamientos 1-4: Ball

**Requests (x4):**
```http
POST /games/a1b2c3d4-e5f6-7890-abcd-ef1234567890/pitches
{"pitchResult": "Ball"}
```

**Responses (primeros 3):**
```json
{ "events": ["PitchThrown"], "atBat": { "result": null } }
```

**Response (4to lanzamiento):**
```json
{
  "game": { ... },
  "atBat": { "result": "Base on Balls", "runsScored": 0, "outsRecorded": 0 },
  "events": ["PitchThrown", "AtBatCompleted"]
}
```

**SSE Emitidos (en orden):**
```
event: PitchThrown (x4)
event: AtBatCompleted
data: {"type":"AtBatCompleted","atBat":{"result":"Base on Balls","outsRecorded":0,...}}
```

> No hay `InningChanged` (0 outs registrados). No hay `ScoreUpdated` (0 carreras).

---

## 3. Hit con 2 Carreras Anotadas

### Contexto
- Inning 5, parte alta, 1 out, corredores en 1ra y 2da
- Cuenta: 2 bolas, 1 strike

### Lanzamiento: In Play (Hit)

**Request:**
```http
POST /games/a1b2c3d4-e5f6-7890-abcd-ef1234567890/pitches
Content-Type: application/json

{
  "pitchResult": "In Play (Hit)",
  "coordinateX": 0.8,
  "coordinateY": 0.1,
  "velocityMPH": 102.3
}
```

**Response (200 OK):**
```json
{
  "game": { "id": "...", "status": "IN_PROGRESS", "homeScore": 3, "awayScore": 4, "currentInning": 5, "isTopInning": true },
  "atBat": { "id": "...", "result": "Hit", "runsScored": 2, "outsRecorded": 0 },
  "events": ["PitchThrown", "AtBatCompleted", "ScoreUpdated"]
}
```

**SSE Emitidos (en orden):**
```
event: PitchThrown
data: {"type":"PitchThrown","pitch":{"pitchResult":"In Play (Hit)","coordinateX":0.8,"coordinateY":0.1,"velocityMPH":102.3,...}}

event: AtBatCompleted
data: {"type":"AtBatCompleted","atBat":{"result":"Hit","runsScored":2,"outsRecorded":0,...}}

event: ScoreUpdated
data: {"type":"ScoreUpdated","game":{"awayScore":4,"homeScore":3,...},"runs_added":2}
```

> **Lógica:** Top inning → suma a `awayScore` (3 → 5, pero ejemplo muestra 4 por simplicidad). `runs_added: 2` refleja las 2 carreras impulsadas.

---

## 4. Foul Ball con 2 Strikes (No Cuenta)

### Contexto
- Cuenta: 1 bola, 2 strikes

**Request:**
```http
POST /games/.../pitches
{"pitchResult": "Foul Ball"}
```

**Response:**
```json
{ "events": ["PitchThrown"], "atBat": { "result": null, "runsScored": 0 } }
```

**SSE:**
```
event: PitchThrown
data: {"type":"PitchThrown","pitch":{"pitchResult":"Foul Ball","ballsBefore":1,"strikesBefore":2,"strikesBefore":2,...}}
```

> El conteo se mantiene en 1 bola, 2 strikes. No hay strikeout.

---

## 5. Out en Juego con Cambio de Inning (3er Out)

### Contexto
- Inning 7, parte baja, 2 outs

**Request:**
```http
POST /games/.../pitches
{"pitchResult": "In Play (Out)"}
```

**Response:**
```json
{
  "game": { "currentInning": 8, "isTopInning": true, "homeScore": 4, "awayScore": 3 },
  "atBat": { "result": "Out", "outsRecorded": 1 },
  "events": ["PitchThrown", "AtBatCompleted", "InningChanged"]
}
```

**SSE Emitidos:**
```
event: PitchThrown
event: AtBatCompleted (result: "Out", outsRecorded: 1)
event: InningChanged
data: {"type":"InningChanged","game":{"currentInning":8,"isTopInning":true,...}}
```

> Transición: `inning 7, bottom` → `inning 8, top`

---

## 6. Secuencia Completa: At-Bat Múltiple con ScoreRuns Manual

### Escenario: Anotación manual de carrera (ej. error defensivo, balk)

**Endpoint interno (no REST público):**
```go
scorekeeperService.ScoreRuns(ctx, gameID, 1)
```

**SSE Emitido:**
```
event: ScoreUpdated
data: {"type":"ScoreUpdated","game":{"homeScore":4,"awayScore":3,...},"runs_added":1}
```

> No hay `PitchThrown` ni `AtBatCompleted` porque no hubo lanzamiento.

---

## Resumen de Eventos por Tipo de Jugada

| Jugada Final | PitchThrown | AtBatCompleted | InningChanged | ScoreUpdated |
|--------------|-------------|----------------|---------------|--------------|
| Ball (1-3) | ✅ | ❌ | ❌ | ❌ |
| Ball (4to → Walk) | ✅ | ✅ | Solo si 3er out | ❌ |
| Strike (1-2) | ✅ | ❌ | ❌ | ❌ |
| Strike (3ro → K) | ✅ | ✅ | Solo si 3er out | ❌ |
| Foul (0-1 strikes) | ✅ | ❌ | ❌ | ❌ |
| Foul (2 strikes) | ✅ | ❌ | ❌ | ❌ |
| In Play (Out) | ✅ | ✅ | Solo si 3er out | ❌ |
| In Play (Hit, 0 carreras) | ✅ | ✅ | ❌ | ❌ |
| In Play (Hit, N carreras) | ✅ | ✅ | ❌ | ✅ |
| ScoreRuns manual | ❌ | ❌ | ❌ | ✅ |

---

## Referencias

- [Endpoints REST →](rest-endpoints.md)
- [Eventos SSE →](sse-events.md)
- [Reglas de negocio →](business-rules.md)
- [Tipos TypeScript →](typescript-types.md)