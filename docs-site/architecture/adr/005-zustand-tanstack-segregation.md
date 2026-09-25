# ADR-005: Separación Zustand (Live/SSE) / TanStack Query (REST)

## Contexto

El frontend consume dos tipos de datos con características opuestas:

| Característica | Datos Live (SSE) | Datos Históricos (REST) |
|----------------|------------------|-------------------------|
| **Fuente** | Push server→client | Pull client→server |
| **Frecuencia** | Alta (cada pitch) | Baja (navegación usuario) |
| **Estado** | Mutaciones incrementales | Reemplazo completo |
| **Caché** | No (efímero) | Sí (stale-while-revalidate) |
| **Suscripción** | Activa mientras vista abierta | Bajo demanda |

Mezclar ambos en una sola librería causa: memory leaks (SSE nunca "stale"), re-fetches innecesarios, complejidad de invalidación.

## Decisión

**Separación estricta de responsabilidades:**

| Librería | Uso Exclusivo | Patrón |
|----------|---------------|--------|
| **Zustand** | Estado SSE tiempo real (juego en vivo, pitches, scoreboard, diamond) | Push → `setState` mutaciones incrementales |
| **TanStack Query** | Datos REST (standings, boxscore, auth, teams, schedule) | Pull → `queryKey` + cache + invalidation |

**Regla dura**: Nunca usar TanStack Query para SSE. Nunca usar Zustand para REST.

## Consecuencias

### Positivas
- **Mental model claro**: "¿Es push o pull?" → decide librería
- **Zustand ligero**: Sin overhead de cache, suscripciones, garbage collection para datos efímeros
- **TanStack Query optimizado**: Deduping, prefetching, background refetch, devtools para datos cacheables
- **Performance**: SSE updates = O(1) state merge; REST = cached/stale-while-revalidate
- **Testing aislado**: Mock store vs mock query client

### Negativas
- **Dos APIs de estado**: Devs deben conocer ambas
- **Puente ocasional**: Ej. SSE `ScoreUpdated` → invalidar `GET /standings` query (hecho via `queryClient.invalidateQueries`)

### Neutrales
- Convención de proyecto, no enforceable por tipos (code review + linting)

## Implementación

```typescript
// store/gameStore.ts (Zustand - SOLO SSE)
export const useGameStore = create<GameState>()(
  subscribeWithSelector((set) => ({
    game: null,
    atBats: [],
    pitches: [],
    // Actions mutan estado incrementalmente
    applyPitchThrown: (pitch) => set((s) => ({ pitches: [...s.pitches, pitch] })),
    applyScoreUpdated: (game) => set({ game }),
  }))
)

// hooks/useStandings.ts (TanStack Query - SOLO REST)
export function useStandings() {
  return useQuery({
    queryKey: ['standings'],
    queryFn: () => api.getStandings(),
    staleTime: 30_000,
  })
}
```

## Alternativas Consideradas

| Opción | Por qué NO |
|--------|------------|
| **Solo TanStack Query** | `setQueryData` para cada pitch = overhead cache, memory leaks, no diseñado para push |
| **Solo Zustand** | Sin cache REST, sin deduping, sin background refetch, sin devtools query |
| **Redux Toolkit** | Overkill, boilerplate, misma dualidad push/pull sin resolver |
| **Signals (Preact/Solid)** | Framework-specific, no portable a Next.js/React |

## Referencias
- [TRD 5.2 Frontend Constraints](../../../docs/TRD.md#52-frontend-constraints-nextjs--react)
- [web/src/lib/utils.ts](../../../web/src/lib/utils.ts) (cn helper)
- [Scorekeeper Panel components](../../../web/src/components/scorekeeper/)