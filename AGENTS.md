# Brandtoon Root Agent Gateway

This file is the **entry point** for AI agents working at repository root (`brandtoon/`).

Its job is to route work correctly between frontend and backend and keep cross-stack changes consistent.

## Scope Routing (critical)

- If task touches only `front/`, follow `front/AGENTS.md`.
- If task touches only `back/`, follow `back/AGENTS.md`.
- If task touches both, treat it as an integration change and follow this file **plus** both scoped guides.

When guidance conflicts, apply this precedence:

1. Direct user instruction
2. Nearest scoped guide (`front/AGENTS.md` or `back/AGENTS.md`)
3. This root gateway

## Repository Map

- `front/` -> React + TypeScript frontend
- `back/` -> Go backend (DDD + CQRS by use case)
- `Justfile` -> project task entrypoints (run commands from repo root unless a task says otherwise)

## Cross-Stack Collaboration Contract

For any feature that spans frontend and backend:

1. Define/confirm backend API contract first (payloads, validation, response, error model).
2. Implement backend behavior while preserving backend architecture rules.
3. Implement frontend integration using the approved frontend state boundaries.
4. Verify end-to-end behavior and edge cases.

Do not introduce frontend assumptions that are not explicitly supported by backend contracts.

## Integration Operating Mode (Full-Stack Default)

Use this mode whenever a task touches both `back/` and `front/`.

### 1) Contract Source of Truth

- Backend API contract is the single source of truth for integration work.
- Define/confirm request/response DTOs, validation, and error model before frontend changes.
- Frontend must consume only documented backend shapes.

### 2) Breaking Change Policy

- If backend contract changes, frontend adaptation is required in the same task, unless explicitly versioned or feature-flagged.
- Do not merge silent contract breaks.

### 3) Full-Stack Definition of Done

An integration task is done only when:

- Backend behavior is implemented under `back/` architecture rules.
- Frontend integration is implemented under `front/` state boundaries.
- Loading, empty, and error states are handled in UI.
- Relevant tests for changed areas pass.
- Contract/documentation is updated when API behavior changes.

### 4) Execution Order (always)

1. Contract
2. Backend implementation
3. Frontend integration
4. End-to-end verification

### 5) Partial Merge Guard

- Avoid partial integration merges that leave frontend and backend incompatible.
- If partial delivery is required, use explicit versioning or feature flags and document rollout constraints.

## Integration PR Checklist (for reviews)

Use this checklist in PR descriptions when the change is **Integration**:

- [ ] Backend contract changes are explicitly documented (request/response/errors/validation).
- [ ] Frontend implementation is aligned with documented backend shapes (no implicit fields).
- [ ] Loading, empty, and error states are implemented in the UI.
- [ ] Relevant tests for touched backend/frontend areas pass.
- [ ] No silent breaking changes, or versioning/feature flag strategy is documented.

## Change Classification

- **Frontend-only**: UI, client state, forms, routing, presentation logic in `front/`
- **Backend-only**: domain/use case/infra/http work in `back/`
- **Integration**: any API contract, DTO, endpoint, or behavior consumed by frontend

Mark the change type before implementation so the correct guide(s) are applied from the start.

## Safety Rules at Root

- Never edit backend from frontend assumptions alone.
- Never edit frontend with undocumented backend response shapes.
- Keep changes minimal and scoped; avoid unrelated refactors.
- If a task is ambiguous (front vs back responsibility), stop and ask for clarification.

## Skill Loading Gateway

When working from root:

- Entering `front/` work -> load skills required by `front/AGENTS.md`.
- Entering `back/` work -> load skills required by `back/AGENTS.md`.
- Full-stack/integration work -> load both sets as needed by touched areas.

This gateway does not replace scoped guides; it enforces coordination between them.

## Quick Start for Agents

Use this checklist before making changes from repository root.

### 1) Frontend-only task example

**Task**: "Add loading skeleton to storyboard page."

- Classify as **Frontend-only**.
- Work only in `front/`.
- Follow `front/AGENTS.md` rules and required frontend skills.
- Do not change backend contracts or assume new API fields.

### 2) Backend-only task example

**Task**: "Add endpoint to regenerate voice track."

- Classify as **Backend-only**.
- Work only in `back/`.
- Follow `back/AGENTS.md` (DDD + CQRS, handler/useCase boundaries, tests).
- If endpoint contract changes, document it clearly for frontend consumers.

### 3) Integration task example

**Task**: "Expose rendering status in backend and show it in frontend timeline."

- Classify as **Integration**.
- Start by defining/confirming backend contract (response shape, errors, validation).
- Implement backend first under `back/` rules.
- Then implement frontend integration under `front/` state boundaries.
- Verify end-to-end behavior and edge cases before considering done.
