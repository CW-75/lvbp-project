---
name: graphify
description: Builds and queries the project knowledge graph at graphify-out/. Use when the user types /graphify, asks to extract, update, or query the graph, map architecture, or run graphify CLI/MCP.
---

# Graphify

Skill portable del repo. Canónica en `.agents/skills/graphify/`. Claude y Windsurf la ven por symlink. La política always-on (consultar el grafo antes de un grep amplio) está en `.agents/rules/graphify.md`, no aquí.

Salida en `graphify-out/` (gitignored). Path por defecto: `.`

## Cuándo usar esta skill vs la regla

- **Regla always-on:** orientarse en código/arquitectura (`query` / `path` / `explain`).
- **Esta skill:** construir o refrescar el grafo, o cuando el usuario dice `/graphify`.

## Pasos

1. Comprobar CLI: `command -v graphify`. Si falta: indicar `pipx install graphifyy` (paquete PyPI `graphifyy`, comando `graphify`). No instalar con `--break-system-packages`. No inventar resultados del grafo.
2. Si el pedido es consulta (`query` / `path` / `explain`): ejecutar el CLI o MCP (`query_graph`, `shortest_path`, `get_node`) y respetar `.agents/rules/graphify.md`.
3. Si el pedido es pipeline (`/graphify`, extract, update):
   - Existe `graphify-out/graph.json` → `graphify update .` (AST, sin API).
   - No existe → `graphify extract .` (código vía AST; docs/imágenes pueden pedir backend LLM).
   - Extra: `--wiki`, `--watch`, `graphify add <url>` solo si el usuario lo pide.
4. Resumir: nodos/aristas o ruta de `GRAPH_REPORT.md`. No volcar `graph.json` ni el reporte entero.

## Consultas

```
graphify query "<pregunta>"
graphify path "<A>" "<B>"
graphify explain "<concepto>"
```

Producto, schema o git: primero `docs/` (PRD, TRD, schema, git_conventions).
