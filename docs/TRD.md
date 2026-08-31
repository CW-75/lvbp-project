# Technical Requirements Document (TRD)

**Proyecto:** Live Baseball GameCast & Standings Engine

## 1. Arquitectura General

Arquitectura híbrida desacoplada:

- **Canal Histórico (REST):** Consultas estáticas (Standings, Boxscore).
- **Canal de Tiempo Real (SSE):** Flujo unidireccional para pitcheos y jugadas.

## 2. Stack Tecnológico

- **Backend:** Go (Golang) >= 1.22.
- **Enrutamiento y Streaming:** `go-chi/chi/v5` (control nativo de `http.Flusher`).
- **Base de Datos Principal:** PostgreSQL 16 (Consultas tipadas con `sqlc` + `pgx/v5`).
- **Caché y Bus de Eventos:** Redis 7 (Pub/Sub con `go-redis/v9`).
- **Frontend (Framework):** Next.js 16+ (App Router, TypeScript 5+).
- **Gestión de Estado (Frontend):** Zustand para eventos SSE (Live) y TanStack Query para REST.
- **Interfaz de Usuario:** Tailwind CSS + shadcn/ui + lucide-react. (Renderizado del diamante y zona de strike mediante SVG nativo).

## 3. Requerimientos No Funcionales (RNF)

- **RNF-1. Latencia:** < 200 ms desde la ingesta del anotador hasta el cliente vía SSE.
- **RNF-2. Concurrencia:** Soportar al menos 5.000 clientes concurrentes por instancia con un consumo < 250 MB RAM.
- **RNF-3. Proxy y Buffering:** Anular buffering intermedio enviando headers `X-Accel-Buffering: no` y `Cache-Control: no-cache`.
