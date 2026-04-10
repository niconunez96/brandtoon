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

## Design System Contract (critical)

All new frontend UI work must inherit the shared Brandtoon visual foundation instead of inventing per-screen styling.

### Canonical tokens

- Font family: **Plus Jakarta Sans**
- Primary action: `#FF6B6B`
- Ink text: `#2D3436`
- Page background: `#fbfcfc`
- Surface ladder:
  - `#f1f2f6` base container
  - `#e5e9eb` high container
  - `#dfe6e9` highest container
- Feedback:
  - error `#d63031`
  - error container `#ffdad6`
  - tertiary container `#00b179`
  - secondary `#006a63`
- Supporting palette samples allowed only when the design needs them:
  - purple `#6B4EE0`
  - pastel pink `#FCE7E7`

### Shape and depth

- Default interactive/card radius: `rounded-2xl` to `rounded-3xl`
- Large hero panels: about `2.5rem` corners
- Use centralized shadow tokens for:
  - overshoot container depth
  - coral sticker/button depth

Do NOT hardcode near-duplicate radii, shadows, or colors inside feature components.

### Responsive usage guidance

- Mobile-first by default.
- `lg`/desktop is the point where the authenticated shell may switch to **fixed sidebar + top bar + content canvas**.
- Tablet should mainly introduce better grouping, not heavier chrome.
- Wide layouts can densify data surfaces, but empty/error/action states must remain readable at narrow widths.
- Reusable guidance lives in `front/docs/design-foundations.md` and must be updated when layout rules change.

## UI Primitive Policy

Before creating feature-specific UI, check whether the need should be composed from shared primitives under `src/shared/components/ui`.

Current foundational primitives include:

- Button variants: primary / secondary / ghost / loading
- Input + textarea
- Badge/chip/tag pills
- Card + section shell
- Slider visual + toggle switch visual
- Toast + action chip
- Topbar / navbar shell
- Sidebar nav
- Data table / asset row
- Metric / progress card

See `front/docs/design-foundations.md` for the explicit token map and foundation rules.

Feature code should compose these primitives first and only extend them when a real product requirement cannot be expressed by the existing layer.

## Layout Rules

- Desktop shells should preserve the **top bar + fixed left sidebar + content canvas** structure when working on the main authenticated SaaS area.
- Mobile shells should preserve the **top app bar + bottom nav + floating action button** pattern unless the flow has a strong product reason to differ.
- Main showcase sections and future documentation should continue mapping UI into these conceptual groups:
  - Hero title
  - Interactions & Buttons
  - Form Architecture
  - Surfaces & Containers
  - Feedback & System Status

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

## Validation Commands

- Canonical frontend validation entrypoint: `just front-check`
- Underlying non-mutating checks:
  - `pnpm --dir front format:check`
  - `pnpm --dir front lint`
  - `pnpm --dir front test`
- CI must stay non-mutating: do NOT use `just front-format` or `pnpm --dir front lint:fix` in validation workflows.

## Styling Rules

1. Tailwind-first with the local project setup only — NEVER use CDN/runtime Tailwind scripts.
2. Load Plus Jakarta Sans from the local dependency, never from remote stylesheet imports.
3. Convert repeated visual values into centralized tokens in `src/styles.css` (or future token files) before reusing them, including custom radii above Tailwind defaults.
4. Plain CSS is allowed only when Tailwind utilities cannot reasonably express the requirement (example: range thumb styling).
5. Keep body copy softer than hero copy; preserve uppercase micro-labels with wide tracking for section metadata.
6. Do not introduce ad-hoc colors or shadows that drift away from the canonical visual contract without updating this file and the frontend baseline skill.

## Decision Log Discipline

When a new frontend decision is made (library, pattern, folder convention), update:

- this `AGENTS.md` if it changes global rules
- `brandtoon-frontend-baseline/SKILL.md` if it changes implementation workflow

Consistency over novelty.
