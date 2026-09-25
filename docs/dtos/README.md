# Documentación de DTOs - Índice

Este directorio contiene la especificación completa de contratos de datos (Data Transfer Objects) para la API REST y Server-Sent Events (SSE) del sistema LVBP.

---

## Estructura

| Archivo | Descripción |
|---------|-------------|
| [`rest-endpoints.md`](rest-endpoints.md) | Endpoints REST: `GET /standings`, `GET /games/{gameId}/boxscore`, `POST /games/{gameId}/pitches` |
| [`sse-events.md`](sse-events.md) | Eventos SSE: `connected`, `PitchThrown`, `AtBatStarted`, `AtBatCompleted`, `InningChanged`, `ScoreUpdated` |
| [`business-rules.md`](business-rules.md) | Reglas de negocio: conteo bolas/strikes, transiciones de inning, validaciones |
| [`flow-examples.md`](flow-examples.md) | Ejemplos de flujos completos: strikeout, walk, hit con carreras |
| [`typescript-types.md`](typescript-types.md) | Referencia TypeScript lista para Frontend (Fases 5-6) |

---

## Convenciones Generales

- **Formato:** JSON (UTF-8)
- **Nombrado:** `camelCase`
- **Fechas:** ISO 8601 (`YYYY-MM-DDTHH:mm:ss.sssZ`)
- **UUIDs:** String estándar
- **Opcionales:** Campo ausente o `null` (omitidos con `omitempty`)

---

## Endpoints Resumen

| Método | Ruta | Request | Response |
|--------|------|---------|----------|
| `GET` | `/standings` | — | `TeamStandings[]` |
| `GET` | `/games/{gameId}/boxscore` | — | `Game` |
| `POST` | `/games/{gameId}/pitches` | `PitchInput` | `PitchResponse` |
| `GET` | `/games/{gameId}/stream` | — | SSE Stream |

---

## Eventos SSE Resumen

| Evento | Descripción | Payload Principal |
|--------|-------------|-------------------|
| `connected` | Conexión establecida | `{status, gameId}` |
| `PitchThrown` | Nuevo lanzamiento | `{type, pitch}` |
| `AtBatStarted` | Inicio turno al bate | `{type, atBat}` |
| `AtBatCompleted` | Turno completado | `{type, atBat}` |
| `InningChanged` | Cambio de inning/mitad | `{type, game}` |
| `ScoreUpdated` | Actualización marcador | `{type, game, runs_added}` |

---

## Versionado

| Versión | Fecha | Cambios |
|---------|-------|---------|
| 1.0.0 | 2025-09-25 | Versión inicial - Contratos completos |