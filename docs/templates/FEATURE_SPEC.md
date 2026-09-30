# Feature Spec: [Feature Name]

**Status:** Draft | In Review | Approved | Implemented
**Date:** YYYY-MM-DD
**Author(s):** [Agent or Human]
**Epic/Requirement:** [Link to issue or PRD section]

---

## 1. Overview
[One-paragraph summary of what this feature is and the business value it provides.]

## 2. Acceptance Criteria
[List of conditions that must be met for the feature to be considered done. Use BDD format if possible (Given/When/Then).]
- AC 1: ...
- AC 2: ...

## 3. Data & Schema Impacts
[List any new tables, columns, or modifications required in PostgreSQL / Redis. Link to `docs/schema.md` if necessary.]
- ...

## 4. Backend & API (Go/REST/SSE)
[Detail new endpoints, SSE channels, or internal core services.]
### Endpoints
- `GET /api/v1/...`
  - Request: `...`
  - Response: `...`

## 5. Frontend & UI (Next.js/React)
[Describe new React components, Zustand state changes, TanStack queries, or UI behavior.]
- ...

## 6. Testing Strategy
[How this feature will be tested. Note if it requires special testcontainers setup or Playwright scenarios.]
- Unit: ...
- E2E: ...

## 7. Known Risks & Constraints
[Any potential bottlenecks, NFR concerns (e.g., latency > 200ms), or dependencies.]
- ...
