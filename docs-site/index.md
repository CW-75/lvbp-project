---
layout: home

hero:
  name: 'LVBP GameCast'
  text: 'Live Baseball GameCast & Standings Engine'
  tagline: 'Documentación técnica completa - Arquitectura, APIs, DTOs y Fases de Implementación'
  image:
    src: /logo-lvbp-light.svg
    alt: LVBP Logo
  actions:
    - theme: brand
      text: 'Comenzar →'
      link: /architecture/system-overview
    - theme: alt
      text: 'Referencia API'
      link: /api-reference
    - theme: alt
      text: 'DTOs Reference'
      link: /dto-reference

features:
  - title: 'Arquitectura Híbrida'
    details: 'REST para consultas históricas (Standings, Boxscore) + SSE para tiempo real (Pitches, Plays, Scoreboard)'
    icon: 🏗️
  - title: 'Motor de Reglas Béisbol'
    details: 'State engine que procesa pitches y calcula transiciones automáticas (balls/strikes/outs/innings)'
    icon: ⚾
  - title: 'Type-Safe End-to-End'
    details: 'Go (sqlc + swaggo) + TypeScript (typedoc) + OpenAPI 3.0 + VitePress para docs unificadas'
    icon: 🔒
  - title: 'Observabilidad'
    details: 'Event Bus Redis Pub/Sub, structured logging, métricas de latencia NFR < 200ms'
    icon: 📊
---

## 📚 Navegación Rápida

| Sección | Descripción |
|---------|-------------|
| [Arquitectura](/architecture/system-overview) | Visión general, diagramas de flujo, ADRs |
| [Fases Implementadas](/phases/01-auth-system) | 8 fases documentadas con código, flujos y decisiones |
| [Referencia DTOs](/dto-reference) | Catálogo completo de DTOs Request/Response/Events |
| [API Reference](/api-reference) | OpenAPI/Swagger UI interactivo |
| [Runbooks](/runbooks/local-dev) | Desarrollo local, despliegue, troubleshooting |

## 🚀 Stack Tecnológico

<div class="tech-stack">
  <span class="badge backend">Go 1.22+</span>
  <span class="badge backend">Chi Router</span>
  <span class="badge backend">sqlc + pgx</span>
  <span class="badge backend">Redis Pub/Sub</span>
  <span class="badge frontend">Next.js 16+</span>
  <span class="badge frontend">TypeScript 5+</span>
  <span class="badge frontend">Zustand + TanStack Query</span>
  <span class="badge frontend">Tailwind + shadcn/ui</span>
  <span class="badge infra">PostgreSQL 16</span>
  <span class="badge infra">Docker Compose</span>
  <span class="badge infra">GitHub Actions</span>
</div>

## 📋 Estado del Proyecto

- ✅ **Fase 1**: Autenticación (JWT + Redis Sessions)
- ✅ **Fase 2**: Modelos de Dominio (Game, AtBat, Pitch, Team, Standings)
- ✅ **Fase 3**: Ingesta de Pitches (REST + Validación)
- ✅ **Fase 4**: Motor Scorekeeper (Reglas béisbol, State Engine)
- ✅ **Fase 5**: Streaming SSE (Tiempo real, Event Bus)
- ✅ **Fase 6**: Standings Handler (REST)
- ✅ **Fase 7**: Boxscore Handler (REST)
- ✅ **Fase 8**: Capa Repositorio (sqlc, PostgreSQL)

<style>
.tech-stack {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin: 1rem 0;
}
.badge {
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
}
.badge.backend { background: #0B1A32; color: white; }
.badge.frontend { background: #007ACC; color: white; }
.badge.infra { background: #336791; color: white; }
</style>