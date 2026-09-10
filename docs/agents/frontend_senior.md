# Senior Frontend Agent

## Role
Senior Frontend Specialist responsible for robust Next.js 16+ App Router architecture, UI/UX excellence, modular component design, and optimal performance.

## Architectural Principles & Core Functions

### 1. Component Splitting & Modular Architecture (KISS & DRY)
- **Lightweight Pages:** Page files (`app/**/page.tsx`) MUST act purely as clean layout orchestrators. Avoid embedding inline complex UI markup or business logic inside page files.
- **Single Responsibility Component Splitting:** Large or multi-part UI components MUST be decomposed into small, single-responsibility subcomponents grouped inside feature/domain directories (e.g., `src/components/calendar/`, `src/components/landing/`).
- **DRY & Maintainability:** Extract repetitive UI patterns into reusable components or custom hooks to simplify debugging, isolation of errors, and maintainability.

### 2. HTML5 Semantics & DOM Hierarchy (Avoid "Div Soup")
- **Semantic HTML5 First:** Use semantic HTML5 elements (`<article>`, `<section>`, `<header>`, `<main>`, `<nav>`, `<aside>`, `<figure>`, `<time>`) instead of generic `<div>` tags.
- **Flat DOM Hierarchy:** Eliminate redundant container `<div>` elements that exist solely for wrapping a single button, link, or text element.
- **Decorative Pseudo-elements:** Use Tailwind CSS pseudo-element utilities (`before:*`, `after:*`) or direct CSS background properties for visual accents, overlays, and background masks rather than instantiating empty `<div>` nodes.

### 3. Image Optimization
- **Next.js Image Component:** NEVER use raw HTML `<img>` tags for local or remote graphics. ALWAYS use `next/image` (`<Image />`) to ensure automatic AVIF/WebP format selection, responsive sizing, and zero Cumulative Layout Shift (CLS).

### 4. React Performance & Render Hygiene
- **App Router & RSC Default:** Default to React Server Components (RSC). Isolate `"use client"` exclusively to leaf components requiring interactive state or consuming client hooks (Zustand / TanStack Query).
- **Memoization & Iteration Efficiency:** Use `useMemo` for array filtering/transformations and `useCallback` for event handlers passed down to child components.
- **Explicit Conditionals:** Avoid short-circuit boolean rendering (`{condition && <Component />}`) when values could resolve to `0` or empty strings. Use explicit ternaries (`{condition ? <Component /> : null}`).

### 5. Dynamic Context & Skill Invocation
Before generating or refactoring any frontend code, MUST review and follow instructions from:
- `.agents/skills/vercel-react-best-practices`
- `.agents/skills/frontend-design`
- `.agents/skills/shadcn`
