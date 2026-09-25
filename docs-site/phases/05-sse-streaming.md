# Fase 5: Streaming SSE (Server-Sent Events)

## Resumen

Canal de tiempo real unidireccional (Server → Client) para actualizaciones live del juego: pitches, at-bats, score, inning changes.

## Arquitectura

```
┌─────────────┐     Redis Pub/Sub      ┌─────────────┐     SSE/HTTP      ┌─────────────┐
│  Scorekeeper │ ─────────────────────▶ │   SSE       │ ────────────────▶ │   Frontend  │
│   Engine     │   game:{id}:events     │  Handler    │  text/event-stream│  (Zustand)  │
└─────────────┘                         └─────────────┘                   └─────────────┘
```

## Endpoint

```
GET /games/{gameId}/stream
Accept: text/event-stream
Authorization: Bearer <token> (opcional para juegos públicos)
```

## Headers Obligatorios (NFR-3)

```go
w.Header().Set("Content-Type", "text/event-stream")
w.Header().Set("Cache-Control", "no-cache")
w.Header().Set("Connection", "keep-alive")
w.Header().Set("X-Accel-Buffering", "no")  // Deshabilita buffering en nginx
```

## Implementación (`internal/handlers/sse/stream_handler.go`)

```go
func (h *StreamHandler) StreamGameEvents(w http.ResponseWriter, r *http.Request) {
    // 1. Headers SSE + anti-buffering
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("X-Accel-Buffering", "no")

    // 2. Verificar Flusher (requerido para streaming)
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
        return
    }

    // 3. Extraer gameId
    gameID := chi.URLParam(r, "gameId")
    if gameID == "" {
        http.Error(w, "gameId is required", http.StatusBadRequest)
        return
    }

    // 4. Suscribirse a canal Redis
    channel := fmt.Sprintf("game:%s:events", gameID)
    pubsub := h.redisClient.GetDB().Subscribe(r.Context(), channel)
    defer pubsub.Close()

    // 5. Confirmar suscripción
    if _, err := pubsub.Receive(r.Context()); err != nil {
        http.Error(w, "Failed to subscribe to events", http.StatusInternalServerError)
        return
    }

    // 6. Evento inicial de conexión
    fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"connected\",\"gameId\":\"%s\"}\n\n", gameID)
    flusher.Flush()

    // 7. Loop de eventos
    ch := pubsub.Channel()
    for {
        select {
        case <-r.Context().Done():
            return  // Cliente desconectado
        case msg, ok := <-ch:
            if !ok {
                return  // Canal cerrado
            }
            // msg.Payload ya es JSON string
            fmt.Fprintf(w, "data: %s\n\n", msg.Payload)
            flusher.Flush()
        }
    }
}
```

## Formato de Eventos SSE

### Conexión Inicial
```
event: connected
data: {"status":"connected","gameId":"uuid"}
```

### Pitch Thrown
```
event: PitchThrown
data: {"type":"PitchThrown","pitch":{"id":"uuid","atBatId":"uuid","pitchNumber":3,"coordinateX":0.1,"coordinateY":0.2,"pitchResult":"Called Strike","ballsBefore":1,"strikesBefore":1,"outsBefore":0,"velocityMPH":92.5,"createdAt":"2026-09-07T19:30:00Z"}}
```

### At-Bat Completado
```
event: AtBatCompleted
data: {"type":"AtBatCompleted","atBat":{"id":"uuid","gameId":"uuid","inning":3,"isTopInning":true,"batterId":"player_123","pitcherId":"player_456","result":"Strikeout","runsScored":0,"outsRecorded":1,"createdAt":"2026-09-07T19:30:00Z"}}
```

### Cambio de Inning
```
event: InningChanged
data: {"type":"InningChanged","game":{"id":"uuid","homeTeamId":"ldc","awayTeamId":"ndm","status":"IN_PROGRESS","homeScore":2,"awayScore":1,"currentInning":4,"isTopInning":true,"startTime":"2026-09-07T19:00:00Z","createdAt":"2026-09-01T10:00:00Z"}}
```

### Actualización de Score
```
event: ScoreUpdated
data: {"type":"ScoreUpdated","game":{...},"runs_added":2}
```

## Cliente Frontend (Zustand + EventSource)

```typescript
// web/src/stores/useGameStore.ts
import { create } from 'zustand'
import { subscribeWithSelector } from 'zustand/middleware'

interface GameEvent {
  type: 'PitchThrown' | 'AtBatStarted' | 'AtBatCompleted' | 'InningChanged' | 'ScoreUpdated'
  pitch?: Pitch
  atBat?: AtBat
  game?: Game
  runs_added?: number
}

interface GameState {
  game: Game | null
  atBats: AtBat[]
  pitches: Pitch[]
  connect: (gameId: string) => void
  disconnect: () => void
  // Actions
  applyPitchThrown: (pitch: Pitch) => void
  applyAtBatCompleted: (atBat: AtBat) => void
  applyInningChanged: (game: Game) => void
  applyScoreUpdated: (game: Game, runs: number) => void
}

export const useGameStore = create<GameState>()(
  subscribeWithSelector((set, get) => ({
    game: null,
    atBats: [],
    pitches: [],
    eventSource: null as EventSource | null,
    
    connect: (gameId: string) => {
      const es = new EventSource(`/api/games/${gameId}/stream`)
      
      es.addEventListener('PitchThrown', (e) => {
        const event = JSON.parse(e.data) as GameEvent
        get().applyPitchThrown(event.pitch!)
      })
      
      es.addEventListener('AtBatCompleted', (e) => {
        const event = JSON.parse(e.data) as GameEvent
        get().applyAtBatCompleted(event.atBat!)
      })
      
      // ... otros listeners
      
      set({ eventSource: es })
    },
    
    disconnect: () => {
      get().eventSource?.close()
      set({ eventSource: null })
    },
    
    applyPitchThrown: (pitch) => set((s) => ({ 
      pitches: [...s.pitches, pitch] 
    })),
    // ...
  }))
)
```

## Reconexión Automática

`EventSource` nativo maneja reconexión:
- Reintento automático ~2-3s tras desconexión
- Al reconectar, cliente debe hacer `GET /games/{id}/boxscore` para sincronizar estado completo
- Opcional: `Last-Event-ID` header para resumir (no implementado aún)

## Escalabilidad

| Aspecto | Solución |
|---------|----------|
| **Multi-instancia** | Redis Pub/Sub fan-out: todas instancias reciben eventos |
| **Memoria** | `http.Flusher` + goroutine por conexión (~2KB c/u) |
| **5,000 clientes** | ~10MB RAM + Redis overhead (cumple NFR-2) |
| **Timeout** | `r.Context().Done()` detecta desconexión cliente |

## Tests

```go
// internal/handlers/sse/stream_handler_test.go
func TestStreamGameEvents_SubscribesAndForwards(t *testing.T) {
    // Setup: mock Redis PubSub
    // Act: GET /games/{id}/stream
    // Assert: Headers correctos, evento connected enviado, mensajes reenviados
}

func TestStreamGameEvents_ClientDisconnect_CleansUp(t *testing.T) {
    // Act: Cliente cierra conexión
    // Assert: pubsub.Close() llamado, goroutine termina
}
```

## Métricas / Observabilidad

- **Conexiones activas**: Gauge `sse_connections_active{game_id}`
- **Eventos enviados**: Counter `sse_events_sent_total{game_id,type}`
- **Latencia publish→deliver**: Histogram `sse_delivery_latency_ms`

## Troubleshooting

| Problema | Causa | Solución |
|----------|-------|----------|
| Eventos no llegan | Proxy buffering (nginx) | `X-Accel-Buffering: no` + `proxy_buffering off` en nginx |
| Conexión cierra a 30s | Load balancer idle timeout | Keep-alive ping cada 15s (TODO) |
| Memoria crece | Goroutines no limpian | Verificar `r.Context().Done()` select case |

## Próximos Pasos

- [ ] Heartbeat/ping cada 15s para evitar timeouts
- [ ] `Last-Event-ID` support para resumir stream
- [ ] Autenticación obligatoria en SSE (actualmente opcional)
- [ ] Compresión gzip para payloads grandes