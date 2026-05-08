# Brandtoon Root Agent Gateway

This file routes work between root, `front/`, and `back/`.

## Canonical Standards

- Repo-wide behavior: `.agents/standards/agent-policy.md`
- Skill authoring: `.agents/standards/skill-authoring.md`
- Runtime skill map: `.atl/skill-registry.md`

## Scope Routing

- `front/` only -> follow `front/AGENTS.md`
- `back/` only -> follow `back/AGENTS.md`
- Both -> treat as **Integration** and follow this file plus both scoped guides

## Repository Map

- `front/` -> React + TypeScript frontend
- `back/` -> Go backend (DDD + CQRS by use case)

## Integration Policy

- Backend API contract is the source of truth.
- Order is always: contract -> backend -> frontend -> verification.
- Do not introduce frontend assumptions not supported by backend contracts.
- If a backend contract changes, adapt the frontend in the same task unless versioned or feature-flagged.
- Do not leave frontend and backend silently incompatible.

## Change Types

- **Frontend-only**: UI, routing, forms, client state in `front/`
- **Backend-only**: domain, use cases, infra, HTTP in `back/`
- **Integration**: shared contract or backend behavior consumed by frontend

## Root Guardrails

- Never edit backend from frontend assumptions alone.
- Never edit frontend against undocumented backend shapes.
- Keep changes minimal and scoped.
- If ownership or policy is unclear, ask the user.
