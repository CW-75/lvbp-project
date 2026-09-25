# Referencia TypeScript para Frontend (Fases 5-6)

Tipos listos para copiar/pegar en el proyecto Frontend (Next.js + TypeScript).

---

## Enums

```typescript
// Estado del juego
type GameStatus =
  | 'SCHEDULED'
  | 'IN_PROGRESS'
  | 'FINAL'
  | 'POSTPONED'
  | 'SUSPENDED';

// Resultado de un lanzamiento (input REST)
type PitchResult =
  | 'Ball'
  | 'Called Strike'
  | 'Swinging Strike'
  | 'Foul Ball'
  | 'In Play (Out)'
  | 'In Play (Hit)';

// Resultado final de un at-bat (output REST/SSE)
type AtBatResult =
  | 'Strikeout'
  | 'Base on Balls'
  | 'Out'
  | 'Hit'
  | null;

// Tipos de eventos SSE
type SSEEventType =
  | 'connected'
  | 'PitchThrown'
  | 'AtBatStarted'
  | 'AtBatCompleted'
  | 'InningChanged'
  | 'ScoreUpdated';
```

---

## DTOs Principales

```typescript
// Tabla de posiciones
interface TeamStandings {
  teamId: number;
  gamesPlayed: number;
  won: number;
  lost: number;
  pct: number;           // 3 decimales (ej: 0.667)
  gamesBehind: number;   // 1 decimal (ej: 1.5)
}

// Juego completo (boxscore)
interface Game {
  id: string;                    // UUID
  homeTeamId: string;            // Código equipo (ej: "LARA")
  awayTeamId: string;
  status: GameStatus;
  homeScore: number;
  awayScore: number;
  currentInning: number;         // 1-9+
  isTopInning: boolean;          // true = visitante batea
  startTime: string;             // ISO 8601
  createdAt: string;             // ISO 8601
}

// Input: Registrar lanzamiento
interface PitchInput {
  pitchResult: PitchResult;      // Obligatorio
  coordinateX?: number;          // Opcional, rango [-1, 1]
  coordinateY?: number;          // Opcional, rango [-1, 1]
  velocityMPH?: number;          // Opcional, > 0
}

// Lanzamiento completo (SSE/DB)
interface Pitch {
  id: string;                    // UUID
  atBatId: string;               // UUID
  pitchNumber: number;           // 1, 2, 3...
  coordinateX: number;
  coordinateY: number;
  pitchResult: PitchResult;
  ballsBefore: number;
  strikesBefore: number;
  outsBefore: number;
  velocityMPH: number | null;
  createdAt: string;             // ISO 8601
}

// Turno al bate
interface AtBat {
  id: string;                    // UUID
  gameId: string;                // UUID
  inning: number;
  isTopInning: boolean;
  batterId: string;
  pitcherId: string;
  result: AtBatResult;           // null = en curso
  runsScored: number;
  outsRecorded: number;          // 0 o 1
  createdAt: string;             // ISO 8601
}

// Response: POST /pitches
interface PitchResponse {
  game: Game;
  atBat: AtBat;
  events: SSEEventType[];        // Eventos generados por esta jugada
}
```

---

## Payloads SSE (Discriminated Unions)

```typescript
// Base común
interface SSEBaseEvent<T extends string, P> {
  type: T;
  payload: P;
}

// connected
interface SSEConnectedPayload {
  status: 'connected';
  gameId: string;
}
type SSEConnectedEvent = SSEBaseEvent<'connected', SSEConnectedPayload>;

// PitchThrown
interface SSEPitchThrownPayload {
  pitch: Pitch;
}
type SSEPitchThrownEvent = SSEBaseEvent<'PitchThrown', SSEPitchThrownPayload>;

// AtBatStarted
interface SSEAtBatStartedPayload {
  atBat: AtBat;
}
type SSEAtBatStartedEvent = SSEBaseEvent<'AtBatStarted', SSEAtBatStartedPayload>;

// AtBatCompleted
interface SSEAtBatCompletedPayload {
  atBat: AtBat;  // result siempre poblado
}
type SSEAtBatCompletedEvent = SSEBaseEvent<'AtBatCompleted', SSEAtBatCompletedPayload>;

// InningChanged
interface SSEInningChangedPayload {
  game: Game;
}
type SSEInningChangedEvent = SSEBaseEvent<'InningChanged', SSEInningChangedPayload>;

// ScoreUpdated
interface SSEScoreUpdatedPayload {
  game: Game;
  runs_added: number;
}
type SSEScoreUpdatedEvent = SSEBaseEvent<'ScoreUpdated', SSEScoreUpdatedPayload>;

// Union de todos los eventos SSE
type SSEEvent =
  | SSEConnectedEvent
  | SSEPitchThrownEvent
  | SSEAtBatStartedEvent
  | SSEAtBatCompletedEvent
  | SSEInningChangedEvent
  | SSEScoreUpdatedEvent;

// Type guard para discriminar por type
function isSSEEvent<T extends SSEEventType>(
  event: SSEEvent,
  type: T
): event is Extract<SSEEvent, { type: T }> {
  return event.type === type;
}
```

---

## Parsing SSE Stream (Helper)

```typescript
// Evento SSE crudo (formato texto)
interface RawSSEMessage {
  event: string;      // ej: "PitchThrown"
  data: string;       // JSON string
}

// Parser genérico
function parseSSEMessage(raw: RawSSEMessage): SSEEvent | null {
  try {
    const payload = JSON.parse(raw.data);
    return { type: raw.event as SSEEventType, payload };
  } catch {
    return null;
  }
}

// Uso con EventSource
const eventSource = new EventSource(`/games/${gameId}/stream`);

eventSource.addEventListener('PitchThrown', (e) => {
  const event = parseSSEMessage({ event: 'PitchThrown', data: e.data });
  if (event && isSSEEvent(event, 'PitchThrown')) {
    console.log('Nuevo lanzamiento:', event.payload.pitch);
    // updatePitchZone(event.payload.pitch);
  }
});

eventSource.addEventListener('AtBatCompleted', (e) => {
  const event = parseSSEMessage({ event: 'AtBatCompleted', data: e.data });
  if (event && isSSEEvent(event, 'AtBatCompleted')) {
    console.log('At-bat finalizado:', event.payload.atBat.result);
    // updateAtBatDisplay(event.payload.atBat);
  }
});

eventSource.addEventListener('ScoreUpdated', (e) => {
  const event = parseSSEMessage({ event: 'ScoreUpdated', data: e.data });
  if (event && isSSEEvent(event, 'ScoreUpdated')) {
    console.log(`+${event.payload.runs_added} carreras`);
    // updateScoreboard(event.payload.game);
  }
});

eventSource.addEventListener('InningChanged', (e) => {
  const event = parseSSEMessage({ event: 'InningChanged', data: e.data });
  if (event && isSSEEvent(event, 'InningChanged')) {
    console.log(`Inning ${event.payload.game.currentInning} - ${event.payload.game.isTopInning ? 'Top' : 'Bottom'}`);
    // updateInningDisplay(event.payload.game);
  }
});
```

---

## API Client Helpers

```typescript
const API_BASE = '/api'; // o import.meta.env.VITE_API_URL

// GET /standings
async function fetchStandings(): Promise<TeamStandings[]> {
  const res = await fetch(`${API_BASE}/standings`);
  if (!res.ok) throw new Error('Failed to fetch standings');
  return res.json();
}

// GET /games/:id/boxscore
async function fetchBoxscore(gameId: string): Promise<Game> {
  const res = await fetch(`${API_BASE}/games/${gameId}/boxscore`);
  if (res.status === 404) throw new Error('Game not found');
  if (!res.ok) throw new Error('Failed to fetch boxscore');
  return res.json();
}

// POST /games/:id/pitches
async function recordPitch(gameId: string, input: PitchInput): Promise<PitchResponse> {
  const res = await fetch(`${API_BASE}/games/${gameId}/pitches`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  });

  const data = await res.json();

  if (!res.ok) {
    const error = new Error(data.error || 'Pitch failed') as Error & { status: number; code?: string };
    error.status = res.status;
    throw error;
  }

  return data;
}

// Helper para validar pitchResult en cliente
const VALID_PITCH_RESULTS: PitchResult[] = [
  'Ball',
  'Called Strike',
  'Swinging Strike',
  'Foul Ball',
  'In Play (Out)',
  'In Play (Hit)',
];

function isValidPitchResult(result: string): result is PitchResult {
  return VALID_PITCH_RESULTS.includes(result as PitchResult);
}
```

---

## Zustand Store Example (Fase 6 - GameCast)

```typescript
import { create } from 'zustand';
import { Game, AtBat, Pitch, SSEEvent } from './types';

interface GameState {
  game: Game | null;
  currentAtBat: AtBat | null;
  pitchHistory: Pitch[];
  events: SSEEvent[];
  
  // Actions
  setGame: (game: Game) => void;
  applyEvent: (event: SSEEvent) => void;
  reset: () => void;
}

export const useGameStore = create<GameState>((set) => ({
  game: null,
  currentAtBat: null,
  pitchHistory: [],
  events: [],

  setGame: (game) => set({ game }),

  applyEvent: (event) =>
    set((state) => {
      const newEvents = [...state.events, event];
      let newGame = state.game;
      let newAtBat = state.currentAtBat;
      let newPitches = state.pitchHistory;

      switch (event.type) {
        case 'PitchThrown':
          newPitches = [...newPitches, event.payload.pitch];
          break;
        case 'AtBatStarted':
          newAtBat = event.payload.atBat;
          break;
        case 'AtBatCompleted':
          newAtBat = event.payload.atBat;
          break;
        case 'InningChanged':
        case 'ScoreUpdated':
          newGame = event.payload.game;
          break;
      }

      return { game: newGame, currentAtBat: newAtBat, pitchHistory: newPitches, events: newEvents };
    }),

  reset: () => set({ game: null, currentAtBat: null, pitchHistory: [], events: [] }),
}));
```

---

## TanStack Query Hooks (Fase 6 - Standings/Boxscore)

```typescript
import { useQuery } from '@tanstack/react-query';
import { fetchStandings, fetchBoxscore, recordPitch } from './api';

// Standings
export function useStandings() {
  return useQuery({
    queryKey: ['standings'],
    queryFn: fetchStandings,
    staleTime: 30_000, // 30s
  });
}

// Boxscore
export function useBoxscore(gameId: string) {
  return useQuery({
    queryKey: ['boxscore', gameId],
    queryFn: () => fetchBoxscore(gameId),
    enabled: !!gameId,
    staleTime: 10_000,
  });
}

// Mutation para registrar pitch
export function useRecordPitch(gameId: string) {
  return useMutation({
    mutationFn: (input: PitchInput) => recordPitch(gameId, input),
    onSuccess: (response) => {
      // Invalidate queries si necesario
      queryClient.invalidateQueries({ queryKey: ['boxscore', gameId] });
    },
  });
}
```

---

## Validaciones Compartidas (Zod)

```typescript
import { z } from 'zod';

export const pitchResultSchema = z.enum([
  'Ball',
  'Called Strike',
  'Swinging Strike',
  'Foul Ball',
  'In Play (Out)',
  'In Play (Hit)',
]);

export const pitchInputSchema = z.object({
  pitchResult: pitchResultSchema,
  coordinateX: z.number().min(-1).max(1).optional(),
  coordinateY: z.number().min(-1).max(1).optional(),
  velocityMPH: z.number().positive().optional(),
});

export type PitchInput = z.infer<typeof pitchInputSchema>;

// Validar en cliente antes de enviar
function validatePitchInput(input: unknown): PitchInput {
  return pitchInputSchema.parse(input);
}
```

---

## Referencias

- [Endpoints REST →](rest-endpoints.md)
- [Eventos SSE →](sse-events.md)
- [Reglas de negocio →](business-rules.md)
- [Ejemplos de flujo →](flow-examples.md)