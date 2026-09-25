# Reglas de Negocio y Validaciones - DTOs

---

## 1. Motor de Conteo (Bolas / Strikes)

| Situación | Resultado | At-bat Completado | Outs |
|-----------|-----------|-------------------|------|
| 4 bolas | **Base on Balls** (Walk) | Sí | 0 |
| 3 strikes | **Strikeout** | Sí | 1 |
| Foul con 0-1 strikes | Strike +1 | No | 0 |
| Foul con 2 strikes | Sin cambio | No | 0 |
| `"In Play (Out)"` | Out en juego | Sí | 1 |
| `"In Play (Hit)"` | Hit en juego | Sí | 0 |

### Lógica de Conteo (Backend: `calculateCount()`)

```go
func calculateCount(pitches []*domain.Pitch) (balls int, strikes int) {
    for _, p := range pitches {
        switch p.PitchResult {
        case "Ball":
            balls++
        case "Called Strike", "Swinging Strike":
            strikes++
        case "Foul Ball":
            if strikes < 2 {
                strikes++
            }
        }
    }
    return balls, strikes
}
```

---

## 2. Transiciones de Inning

| Condición | Acción |
|-----------|--------|
| 3 outs alcanzados | Cambio de mitad de inning |
| Era `isTopInning: true` | `isTopInning = false` (mismo inning) |
| Era `isTopInning: false` | `isTopInning = true` + `currentInning++` |

### Evento Emitido

Siempre que hay transición → `InningChanged` con `Game` actualizado.

---

## 3. Estados de Juego Válidos para Ingesta

El endpoint `POST /games/{gameId}/pitches` **solo acepta** requests si:

```
game.status === "IN_PROGRESS"
```

### Estados Posibles (`GameStatus`)

| Estado | Descripción | Permite Ingesta |
|--------|-------------|-----------------|
| `SCHEDULED` | Programado, no iniciado | ❌ |
| `IN_PROGRESS` | En curso | ✅ |
| `FINAL` | Finalizado | ❌ |
| `POSTPONED` | Pospuesto | ❌ |
| `SUSPENDED` | Suspendido | ❌ |

---

## 4. Validaciones de Entrada (`PitchInput`)

### `pitchResult` - Obligatorio

Debe ser uno de los 6 valores válidos (ver `rest-endpoints.md`). Validación en backend:

```go
func isValidPitchResult(result string) bool {
    validResults := map[string]bool{
        "Ball":             true,
        "Called Strike":    true,
        "Swinging Strike":  true,
        "Foul Ball":        true,
        "In Play (Out)":    true,
        "In Play (Hit)":    true,
    }
    return validResults[result]
}
```

### `coordinateX` / `coordinateY` - Opcionales

- Rango esperado: `[-1.0, 1.0]`
- Representan posición en zona de strike normalizada
- `null` o ausente = no proporcionado

### `velocityMPH` - Opcional

- Valor esperado: `> 0` (típico 70-105)
- `null` o ausente = no medido

---

## 5. Lógica de Anotación de Carreras

### Endpoint Interno: `ScoreRuns(ctx, gameID, runs)`

```go
if game.IsTopInning {
    game.AwayScore += runs  // Top = visitante batea
} else {
    game.HomeScore += runs  // Bottom = local batea
}
```

### Evento `ScoreUpdated`

Se emite **siempre** que `runs > 0`, con:
- `game` actualizado
- `runs_added`: número de carreras en la jugada

---

## 6. Flujo de Eventos por Tipo de Jugada

| Jugada | Eventos SSE Emitidos (orden) | `PitchResponse.events` |
|--------|------------------------------|------------------------|
| Ball (1-3) | `PitchThrown` | `["PitchThrown"]` |
| Ball (4to) | `PitchThrown` → `AtBatCompleted` | `["PitchThrown","AtBatCompleted"]` |
| Strike (1-2) | `PitchThrown` | `["PitchThrown"]` |
| Strike (3ro) | `PitchThrown` → `AtBatCompleted` (+ `InningChanged` si 3er out) | `["PitchThrown","AtBatCompleted",...]` |
| Foul (0-1 strikes) | `PitchThrown` | `["PitchThrown"]` |
| Foul (2 strikes) | `PitchThrown` | `["PitchThrown"]` |
| In Play (Out) | `PitchThrown` → `AtBatCompleted` (+ `InningChanged` si 3er out) | `["PitchThrown","AtBatCompleted",...]` |
| In Play (Hit) + 0 carreras | `PitchThrown` → `AtBatCompleted` | `["PitchThrown","AtBatCompleted"]` |
| In Play (Hit) + N carreras | `PitchThrown` → `AtBatCompleted` → `ScoreUpdated` | `["PitchThrown","AtBatCompleted","ScoreUpdated"]` |

---

## 7. Restricciones de Integridad

| Regla | Validación |
|-------|------------|
| Un juego debe existir | 404 si `gameId` no existe |
| Solo un at-bat activo por juego | 409 si `atBat.Result != nil` |
| Juego debe estar en progreso | 409 si `status != IN_PROGRESS` |
| `gameId` debe ser UUID válido | 400 formato inválido |
| JSON debe ser válido | 400 parse error |

---

## Referencias

- [Endpoints REST →](rest-endpoints.md)
- [Eventos SSE →](sse-events.md)
- [Ejemplos de flujo →](flow-examples.md)
- [Tipos TypeScript →](typescript-types.md)