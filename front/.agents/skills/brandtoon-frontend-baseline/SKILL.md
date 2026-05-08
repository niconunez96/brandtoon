---
name: brandtoon-frontend-baseline
description: Standard baseline for Brandtoon frontend implementation using Vite, React, TypeScript, Tailwind, TanStack Query, Zustand, React Hook Form, and Zod. Use for any feature work in front/.
---

# Brandtoon Frontend Baseline

## Purpose

Provide a **simple, repeatable, low-complexity** frontend workflow for contributors with mixed frontend experience.

## Canonical Ownership

- Frontend policy lives in `front/AGENTS.md`.
- Design-system and layout policy live in `front/AGENTS.md` and `front/docs/design-foundations.md`.
- This skill keeps the **operational workflow** for implementing work in `front/`.

## Core Principles

1. **Simplicity first**: choose the clearest implementation.
2. **Separation of concerns**: server state != UI state != form state.
3. **Backend-centric domain logic**: frontend orchestrates UX, not heavy processing.
4. **Consistency beats cleverness**: use the same patterns everywhere.

## State Decision Matrix

Use this table before writing state code:

- API response/cache/retry/invalidation/polling -> **TanStack Query**
- Modal open/close, tabs, filters, wizard step, ephemeral selections -> **Zustand**
- Form fields, validation, submit lifecycle -> **RHF + Zod**

If uncertain, default to TanStack Query for anything coming from backend.

## Working Shape

- Prefer feature-first structure.
- Keep route pages thin and push reusable UI/logic downward.
- Add store folders only when local UI orchestration actually exists.
- Do not introduce new layers unless the pattern repeats.

## Feature Implementation Checklist

For each new feature:

1. Define domain types and validation schema first (TypeScript + Zod).
2. Implement API layer functions in `feature/api`.
3. Add Query hooks for read/mutation flows.
4. Add Zustand store only for local UI orchestration (if needed).
5. Build UI with small components and Tailwind utilities.
6. Handle loading/error/empty/success states explicitly.

## Validation Commands

- Required local implementation loop for frontend feature/fix work:
  1. Run `just front-format` after adding or changing code so Biome can fix imports and formatting.
  2. Run `just front-check` as the canonical feedback loop.
  3. Iterate on failures until `just front-check` passes.
- Canonical frontend validation entrypoint: `just front-check`
- Underlying non-mutating commands:
  - `pnpm --dir front format:check`
  - `pnpm --dir front lint`
  - `pnpm --dir front test`
- CI and review workflows must stay non-mutating: never replace validation with `just front-format` or `pnpm --dir front lint:fix`.

## Frontend Test Baseline

- Default to **unit tests and focused component tests** for frontend behavior.
- Prefer mocking hooks, navigation, and API/service seams over MSW-driven app-shell flows when the goal is page or component behavior.
- Keep `App`-level tests thin and limited to routing composition or other wiring that cannot be validated lower in the tree.
- Treat **fast feedback** as a requirement: aim for **under 500ms per test/spec where practical**.
- If a test becomes slow because it crosses too many boundaries, move that coverage down to the owning component/page.
- Every test must earn its place by increasing **real confidence** in behavior, contracts, decisions, or failure handling.
- Do **not** preserve tests for coverage theater. If a test only proves superficial rendering and does not meaningfully protect the system, delete it.
- Prefer fewer high-value tests over broad low-signal suites that slow development while providing fake safety.

## Async Job UX Pattern (for render/generation flows)

For backend-driven long jobs:

1. Create job via mutation.
2. Store returned `jobId`.
3. Poll job status with Query (or use SSE/WebSocket when available).
4. Render status timeline: `queued -> processing -> done/failed`.
5. On completion, invalidate related queries and display result.

Keep this logic in feature hooks/components; avoid app-wide complexity.

## Guardrails

- No premature abstractions.
- No custom state framework wrappers.
- No mixing backend orchestration logic into generic UI components.
- No hidden side effects in stores.

If a design-system or layout rule is needed, reference `front/AGENTS.md` and `front/docs/design-foundations.md` instead of restating it here.

The best frontend here is the one a new teammate can understand in minutes.
