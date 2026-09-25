# Architecture Decision Records (ADRs)

Registro de decisiones arquitectónicas significativas para el proyecto LVBP GameCast.

## Lista de ADRs

| # | Título | Estado | Fecha |
|---|--------|--------|-------|
| [001](./001-hybrid-rest-sse.md) | Arquitectura Híbrida REST + SSE | Aceptado | 2026-09-24 |
| [002](./002-sqlc-over-orm.md) | sqlc en lugar de ORM | Aceptado | 2026-09-24 |
| [003](./003-scorekeeper-explicit-input.md) | Scorekeeper sigue input explícito | Aceptado | 2026-09-24 |
| [004](./004-redis-pubsub-eventbus.md) | Redis Pub/Sub como Event Bus | Aceptado | 2026-09-24 |
| [005](./005-zustand-tanstack-segregation.md) | Separación Zustand / TanStack Query | Aceptado | 2026-09-24 |

## Formato ADR

Cada ADR sigue el formato:

```markdown
# ADR-XXX: Título

## Contexto
¿Cuál es el problema que estamos resolviendo?

## Decisión
¿Qué decidimos hacer?

## Consecuencias
### Positivas
- ...

### Negativas
- ...

### Neutrales
- ...
```

## Cómo Añadir un Nuevo ADR

1. Crear archivo `NNN-titulo-kebab-case.md` en esta carpeta
2. Usar el siguiente número secuencial
3. Actualizar esta tabla de contenidos
4. Referenciar desde docs relevantes