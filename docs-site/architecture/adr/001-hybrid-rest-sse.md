# ADR-001: Arquitectura Híbrida REST + SSE

## Contexto

El sistema necesita servir dos tipos de clientes con requisitos muy diferentes:

1. **Consultas Históricas/Estáticas**: Standings, Boxscore, Autenticación - acceso esporádico, cacheable, request-response clásico
2. **Tiempo Real**: Pitches, plays, scoreboard live - latencia crítica (<200ms), flujo unidireccional servidor→cliente, miles de conexiones concurrentes

Usar solo REST con polling sería ineficiente para tiempo real. Usar solo WebSockets añade complejidad innecesaria para consultas estáticas. WebSockets también tienen problemas con proxies corporativos, load balancers y reconexión.

## Decisión

**Arquitectura híbrida de dos canales:**

| Canal | Protocolo | Endpoints | Patrón |
|-------|-----------|-----------|--------|
| **Histórico** | REST/HTTPS | `GET /standings`, `GET /games/{id}/boxscore`, `POST /auth/*` | Request-Response |
| **Tiempo Real** | SSE (Server-Sent Events) | `GET /games/{id}/stream` | Pub/Sub Unidireccional |

**Implementación SSE:**
- Go nativo `http.Flusher` via `chi` router
- Headers obligatorios: `Content-Type: text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive`, `X-Accel-Buffering: no`
- Redis Pub/Sub como backbone: `game:{gameId}:events`
- Reconexión automática nativa en navegador (EventSource API)

## Consecuencias

### Positivas
- **Separación de responsabilidades clara**: REST para CRUD, SSE para streaming
- **Simplicidad operativa**: SSE usa HTTP/1.1 o HTTP/2 estándar, pasa firewalls/proxies sin config especial
- **Reconexión automática**: `EventSource` del navegador maneja reconexión nativamente
- **Escalabilidad**: Redis Pub/Sub fan-out permite múltiples instancias backend
- **Debugging trivial**: `curl -N <url>` muestra stream en terminal
- **NFR-3 cumplido**: Headers anti-buffering nativos

### Negativas
- **Unidireccional**: Cliente no puede enviar datos por SSE (usa REST separado para eso)
- **Límite conexiones HTTP/1.1**: 6 conexiones por host (mitigado con HTTP/2)
- **No binario**: Solo texto/JSON (aceptable para payloads béisbol)

### Neutrales
- Requiere sticky sessions o Redis compartido para SSE en multi-instancia
- Frontend usa dos librerías estado: Zustand (SSE live) + TanStack Query (REST)

## Alternativas Consideradas

| Opción | Por qué NO |
|--------|------------|
| **Solo REST + Polling** | Latencia alta, overhead HTTP, no cumple NFR-1 |
| **Solo WebSockets** | Overengineering para consultas estáticas, complejidad reconexión, proxies |
| **gRPC + gRPC-Web** | Complejidad tooling, browser support limitado, overkill |
| **GraphQL Subscriptions** | Complejidad schema, WebSocket transport, YAGNI |

## Referencias
- [NFR-1 Latencia < 200ms](../system-overview.md#nfrs-cumplidos)
- [NFR-3 Proxy Buffering](../system-overview.md#nfrs-cumplidos)
- [TRD Sección 1](../architecture/system-overview.md)