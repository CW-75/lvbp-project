---
trigger: always_on
description: Consulta el grafo Graphify en graphify-out/ para preguntas de código y arquitectura.
---

# Graphify (fuente portable)

Política única de consulta. Cursor: `.cursor/rules/graphify.mdc`. Windsurf: `.windsurf/rules/graphify.md`. Pipeline `/graphify`: skill `.agents/skills/graphify/SKILL.md`. No duplicar el cuerpo.

El grafo vive en `graphify-out/` (gitignored).

## Prioridad

- Producto, tech, schema o git: empieza por `docs/` (PRD, TRD, schema, git_conventions), como en `AGENTS.md`.
- Código o arquitectura: si existe `graphify-out/graph.json`, consulta el grafo **antes** de un grep/glob amplio.

## Cómo consultar

- Preferir `graphify query "<pregunta>"` (CLI) o `query_graph` (MCP).
- Relaciones: `graphify path "<A>" "<B>"` / `shortest_path`.
- Conceptos: `graphify explain "<concepto>"` / `get_node`.
- Si existe `graphify-out/wiki/index.md`, navegarlo en lugar de volcar el grafo crudo.
- Leer `graphify-out/GRAPH_REPORT.md` solo en revisión amplia de arquitectura, o si query/path/explain no bastan. No cargarlo por defecto (es grande).

## Fallbacks

- Si no hay CLI/MCP `graphify`, usar `docs/` y lecturas puntuales. No inventar resultados del grafo.
- Tras cambios de código sustanciales, ejecutar `graphify update .` solo si el CLI está disponible (AST, sin costo de API). Si no está, omitir en silencio.
