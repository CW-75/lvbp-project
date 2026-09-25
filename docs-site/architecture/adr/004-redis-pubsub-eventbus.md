# ADR-004: Redis Pub/Sub como Event Bus

## Contexto

El Scorekeeper Engine necesita notificar a múltiples consumidores (SSE handlers, futuros webhooks, analytics) cuando ocurren eventos de dominio:
- `PitchThrown`
- `AtBatStarted` / `AtBatCompleted`
- `InningChanged`
- `ScoreUpdated`

Opciones: Canal Go interno, RabbitMQ/Kafka, Redis Pub/Sub, NATS.

## Decisión

**Usar Redis Pub/Sub como Event Bus interno.**

- Canal por juego: `game:{gameId}:events`
- Payload: JSON con `{type, ...data}`
- Publicación fire-and-forget (`_ = eventBus.PublishEvent(...)`)
- Suscripción en SSE Handler via `redis.Subscribe(ctx, channel)`

## Consecuencias

### Positivas
- **Desacoplamiento total**: Engine no conoce consumidores
- **Escalabilidad horizontal**: Múltiples instancias API comparten mismo Redis
- **Simplicidad operativa**: Ya usamos Redis para sessions/cache
- **Fan-out nativo**: Un publish → N suscriptores SSE
- **Bajo overhead**: Pub/Sub en memoria, sub-millisecond latency
- **Replay opcional**: Puede migrar a Redis Streams si se necesita persistencia

### Negativas
- **No persistente**: Si no hay suscriptores, evento se pierde (aceptable: SSE reconecta y pide estado via REST)
- **Sin ordering global**: Solo ordering por canal (suficiente: un canal por juego)
- **At-least-once**: Posibles duplicados en redelivery (idempotencia en consumidor)

### Neutrales
- Requiere Redis cluster para HA en producción
- Payload JSON → parsing en consumidor (rápido en Go/TS)

## Alternativas Consideradas

| Opción | Por qué NO |
|--------|------------|
| **Canal Go interno (chan)** | No funciona multi-instancia, acopla Engine a handlers |
| **RabbitMQ** | Overhead operacional, complejidad routing, YAGNI |
| **Kafka** | Overkill para < 10k events/sec, operacionalmente pesado |
| **NATS** | Otra dependencia, Redis ya está en stack |
| **PostgreSQL LISTEN/NOTIFY** | Limitado a una conexión, no escala bien para muchos suscriptores |

## Referencias
- [EventPublisher interface](../../../backend/internal/core/ports/events.go)
- [Redis PubSub impl](../../../backend/internal/infrastructure/pubsub/redis.go)
- [SSE StreamHandler](../../../backend/internal/handlers/sse/stream_handler.go)
- [Scorekeeper Service](../../../backend/internal/core/services/scorekeeper.go)