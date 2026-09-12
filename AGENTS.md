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

## Context Optimization Rules
- **Telegraphic Style:** Always use concise, telegraphic language in agent instructions and internal reasoning. Omit conversational filler.
- **Abbreviations:** Use established acronyms (e.g., LBGC, NFR, TRD, PRD) consistently to save tokens.
- **Agent Consolidation:** Prefer a single agent for highly overlapping domains (e.g., Backend & Data) to minimize context switching overhead.
- **Tool Protocol:** Use native tools (`view_file`, `replace_file_content`, `run_command`) directly instead of outputting code blocks in chat. Avoid `cat`/`grep` in bash.
