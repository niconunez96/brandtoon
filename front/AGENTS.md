# Brandtoon Frontend Agent Guide

This document defines the **default frontend baseline** for this project.

## Mission

Build a frontend that is:

- **Simple to understand** for non-frontend experts
- **Consistent** across features and contributors
- **Fast** in development and runtime
- **Boring in architecture** (intentionally): no accidental complexity

## Approved Stack (Default)

- Build tool: **Vite**
- UI: **React + TypeScript**
- Styling: **Tailwind CSS**
- Routing: **React Router**
- Server state: **TanStack Query**
- Client/UI state: **Zustand**
- Forms: **React Hook Form + Zod**
- Package manager: **pnpm**

## Non-Goals (avoid these unless explicitly requested)

- No Next.js
- No Redux unless a concrete scaling need appears
- No global event bus pattern
- No over-engineered abstractions for small features

## State Management Policy (critical)

Keep state split by responsibility:

1. **Server state** (API data, async jobs, polling, cache) -> TanStack Query
2. **Client/UI state** (modals, filters, stepper state, local selections) -> Zustand
3. **Form state** (input, validation, submission) -> React Hook Form + Zod

Do not mix these concerns in one store.

## Simplicity Rules

1. Prefer feature folders over deep architectural layering.
2. Start with straightforward implementations; abstract only after duplication is proven.
3. Co-locate code near usage (feature-first structure).
4. Keep components focused and small.
5. Avoid clever patterns that reduce readability.

## Frontend Boundaries

The backend owns heavy processing (image generation, animation rendering, voice pipeline).
The frontend should only:

- Collect user input
- Trigger jobs
- Show progress/state transitions
- Display outputs and errors clearly

No heavy business logic in the frontend unless strictly required.

## Skill Loading Rules

When working in `front/`, load these skills:

- `brandtoon-frontend-baseline` for architecture, state boundaries, and folder conventions
- `vercel-react-best-practices` for performance-sensitive React decisions
- `web-design-guidelines` for UI review and accessibility checks

Use `frontend-design` only when explicitly asked for visual redesign/creative UI direction.

## Decision Log Discipline

When a new frontend decision is made (library, pattern, folder convention), update:

- this `AGENTS.md` if it changes global rules
- `brandtoon-frontend-baseline/SKILL.md` if it changes implementation workflow

Consistency over novelty.
