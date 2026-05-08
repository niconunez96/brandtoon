# Brandtoon Frontend Policy

Frontend agents must also follow `../.agents/standards/agent-policy.md`.

## Stack Baseline

- Vite
- React + TypeScript
- Tailwind CSS
- React Router
- TanStack Query for server state
- Zustand for client/UI state
- React Hook Form + Zod for forms

## Frontend Architecture

- Prefer feature-first structure and straightforward components.
- Do not mix server state, UI state, and form state in one store.
- Keep heavy business logic in the backend unless the UI strictly requires it.

## Design System Policy

- Reuse the shared Brandtoon visual foundation instead of inventing per-screen styling.
- Compose shared primitives under `src/shared/components/ui` before adding feature-specific UI.
- Keep canonical tokens and layout guidance aligned with `front/docs/design-foundations.md`.
- Do not introduce drifting colors, shadows, or radii without updating the shared foundation.

## Product Shell Rules

- Desktop product surfaces default to **top bar + fixed sidebar + content canvas**.
- Mobile product surfaces default to **top app bar + bottom nav + floating action button**.
- Loading, empty, and error states must remain readable across narrow widths.

## Frontend Guardrails

- No Next.js unless explicitly requested.
- No Redux unless a concrete scaling need is established.
- No global event bus.
- If a new frontend convention is needed, ask the user instead of inventing it.
