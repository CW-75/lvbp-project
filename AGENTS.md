# AI Agent Context

Welcome, AI Agent. To understand the full context of this project, please refer to the following documents located in the `./docs` directory:

- [Product Requirements Document (PRD)](./docs/PRD.md): Contains the product vision, goals, and features.
- [Technical Requirements Document (TRD)](./docs/TRD.md): Contains the technical architecture, stack, and constraints.
- [Database Schema](./docs/schema.md): Full DDL schema for DB migrations and sqlc queries.
- [Database Schema Summary](./docs/schema_summary.md): Lightweight entity overview for high-level tasks.
- [Git Conventions](./docs/git_conventions.md): Strict guidelines on branching, commits, and environment reuse.

**CRITICAL REQUIREMENT:** You MUST read and follow the [Git Conventions](./docs/git_conventions.md) before working on ANY requirement or change.

Please read these documents before starting your tasks to ensure alignment with the project's objectives and guidelines.

## Communication Guidelines

- **IMPORTANT**: Even if the user prompts in English, you MUST bring all solutions, suggestions, and explanations translated in Spanish in the Agent chat.
- **IMPORTANT**: Any proposed implementation plan, task list, walkthrough, or other generated markdown artifact MUST also be written entirely in Spanish.

## Available Core Skills

To reduce context consumption, load the following skills **ONLY** when your task requires it:
- **UI / Frontend:** Load the `shadcn` skill if you work with UI components, and `frontend-design` if you design new interfaces. Use `vercel-react-best-practices` for optimization.
- **Database / Backend:** Load the `supabase-postgres-best-practices` skill BEFORE modifying schemas, queries, or migrations.
- **Architecture:** Load `graphify` to query the project's knowledge graph.

## Context Optimization Rules
- **Telegraphic Style:** Always use concise, telegraphic language in agent instructions and internal reasoning. Omit conversational filler.
- **Abbreviations:** Use established acronyms (e.g., LBGC, NFR, TRD, PRD) consistently to save tokens.
- **Agent Consolidation:** Prefer a single agent for highly overlapping domains (e.g., Backend & Data) to minimize context switching overhead.
- **Tool Protocol:** Use native tools (`view_file`, `replace_file_content`, `run_command`) directly instead of outputting code blocks in chat. Avoid `cat`/`grep` in bash.
- **Git Flow:** When starting a NEW requirement, always branch off from `develop`.

## SDD & Knowledge Persistence
- **Spec-Driven Development:** Before creating an execution plan for a new feature, you MUST instantiate a feature spec using `docs/templates/FEATURE_SPEC.md` and save it to `docs/specs/`. No code should be planned without a spec.
- **Memory Persistence:** Before debugging a complex issue, read `docs/MEMORY.md`. If you resolve a difficult bug or discover a new project constraint, you MUST append the lesson learned to `docs/MEMORY.md`.
- **ADRs:** If you make a major architectural decision, generate an Architecture Decision Record in `docs-site/architecture/adr/`.
