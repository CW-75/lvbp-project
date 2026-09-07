# Git Conventions & Guidelines

This document outlines the Git conventions that must be strictly followed when working on this project.

## General Principles

1. **Reusing Branches/Worktrees for Consecutive Changes:**
   - **CRITICAL:** Do NOT create other branches or worktrees if the same files are modified by a requirement on chat. Reuse the current branch or worktree to avoid fragmenting the workspace unnecessarily and creating conflicts across multiple branches.
   
2. **Commit Messages:**
   - Use descriptive commit messages following the conventional commits standard (e.g., `feat:`, `fix:`, `chore:`, `docs:`, `refactor:`).
   - Keep the summary line concise (under 50-70 characters).

3. **Branch Naming (When a new branch is explicitly needed):**
   - Feature: `feature/<feature-name>`
   - Bugfix: `hotfix/<bug-name>`
   - Chore/Docs: `chore/<chore-name>`

4. **Git Flow & Workspace Management:**
   - **Lazy Evaluation:** DO NOT invoke read-only tools (`git_status`, `git_log`) unless you lack the minimum context required to make a decision.
   - **Single Responsibility Principle (SRP) for Environments:** Each working directory must have a single purpose. Never mix `feature/` development with `hotfix/` resolution in the same directory.
   - **Workflow Decision Matrix:** Before executing any Git modification tool, evaluate the user's request.
     - **Scenario A (Linear / Sequential Work):** If continuing an ongoing ticket, creating a standard new `feature/`, or syncing `develop`, use native git commands or terminal in the main directory.
     - **Scenario B (Interruption / Parallel Work):** If handling an urgent bug (`hotfix/`), reviewing a third-party Pull Request, or starting a blocking task without losing unsaved changes, initiate the parallel flow. Use `git worktree add` (e.g., to a `../hotfix-temp` directory). Execute changes in that isolated directory. Upon completion, use `git worktree remove`.
   - **Strict Constraint:** If unsure whether the task requires isolation, assume Scenario A (KISS) to avoid environment over-engineering, unless uncommitted changes are at risk.
   - **Agent Orchestration:** Every worker should work on the same branch or worktree depending on the context assigned by the orchestrator.
