---
name: brandtoon-frontend-baseline
description: Standard baseline for Brandtoon frontend implementation using Vite, React, TypeScript, Tailwind, TanStack Query, Zustand, React Hook Form, and Zod. Use for any feature work in front/.
---

# Brandtoon Frontend Baseline

## Purpose

Provide a **simple, repeatable, low-complexity** frontend workflow for contributors with mixed frontend experience.

## Core Principles

1. **Simplicity first**: choose the clearest implementation.
2. **Separation of concerns**: server state != UI state != form state.
3. **Backend-centric domain logic**: frontend orchestrates UX, not heavy processing.
4. **Consistency beats cleverness**: use the same patterns everywhere.

## Standard Stack

- Vite + React + TypeScript
- Tailwind CSS
- React Router
- TanStack Query (server state)
- Zustand (client/UI state)
- React Hook Form + Zod (forms)
- pnpm

## State Decision Matrix

Use this table before writing state code:

- API response/cache/retry/invalidation/polling -> **TanStack Query**
- Modal open/close, tabs, filters, wizard step, ephemeral selections -> **Zustand**
- Form fields, validation, submit lifecycle -> **RHF + Zod**

If uncertain, default to TanStack Query for anything coming from backend.

## Suggested Folder Convention

```txt
src/
  app/                  # app bootstrapping, providers, router setup
  routes/               # route-level pages
  features/
    auth/
      api/              # query/mutation functions
      components/
      hooks/
      schemas/          # zod schemas
      store/            # zustand store (only if needed)
      types/
    generation/
      api/
      components/
      hooks/
      schemas/
      store/
      types/
  shared/
    components/         # reusable presentational pieces
    lib/                # query client, utils, config
    styles/
```

Do not introduce more layers unless there is a clear, repeated need.

## Feature Implementation Checklist

For each new feature:

1. Define domain types and validation schema first (TypeScript + Zod).
2. Implement API layer functions in `feature/api`.
3. Add Query hooks for read/mutation flows.
4. Add Zustand store only for local UI orchestration (if needed).
5. Build UI with small components and Tailwind utilities.
6. Handle loading/error/empty/success states explicitly.

## Validation Commands

- Canonical frontend validation entrypoint: `just front-check`
- Underlying non-mutating commands:
  - `pnpm --dir front format:check`
  - `pnpm --dir front lint`
  - `pnpm --dir front test`
- CI and review workflows must stay non-mutating: never replace validation with `just front-format` or `pnpm --dir front lint:fix`.

## Async Job UX Pattern (for render/generation flows)

For backend-driven long jobs:

1. Create job via mutation.
2. Store returned `jobId`.
3. Poll job status with Query (or use SSE/WebSocket when available).
4. Render status timeline: `queued -> processing -> done/failed`.
5. On completion, invalidate related queries and display result.

Keep this logic in feature hooks/components; avoid app-wide complexity.

## Tailwind Usage Rules

1. Prefer utility classes directly in component markup.
2. Extract component only when repeated or too large.
3. Keep design tokens centralized when introduced (colors/spacing).
4. Avoid custom CSS unless utilities cannot express the intent.

## Guardrails

- No premature abstractions.
- No custom state framework wrappers.
- No mixing backend orchestration logic into generic UI components.
- No hidden side effects in stores.

The best frontend here is the one a new teammate can understand in minutes.
