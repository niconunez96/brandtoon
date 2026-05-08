# Brandtoon Backend Policy

Backend agents must also follow `../.agents/standards/agent-policy.md`.

## Architecture Contract

- Preserve strict DDD + CQRS by use case.
- Dependency flow is `main.go -> routes.go -> handlers -> useCases -> domain`.
- Domain never imports transport or infrastructure.
- Shared dependency initialization lives in `bounded_contexts/shared/diContainer.go`.

## Backend Structure

- Base path: `bounded_contexts/{context}/{aggregate}/`
- Aggregate layers: `domain/`, `useCases/`, `infra/`
- Shared modules live under `bounded_contexts/shared/`
- Aggregate package names stay explicit: `{aggregate}domain`, `{aggregate}usecases`, `{aggregate}repo`, `{aggregate}http`

## Application Rules

- One use case per file.
- Use cases and handlers are functions, not service structs.
- Commands do side effects; queries do reads.
- Use-case DTOs are the canonical payload types when transport shape matches.
- Aggregate-owned endpoints, handlers, and use cases stay in the same aggregate.
- Domain and use-case naming must stay provider-agnostic.

## Infrastructure Rules

- Huma is the standard HTTP framework.
- HTTP routes are grouped in `infra/http/routes.go`.
- Repository adapters live in `infra/repo/`.
- Third-party clients must be initialized through the shared DI container.

## Testing Policy

- TDD is mandatory for backend changes in `useCases` and `infra/http`.
- Tests for those areas are required when touched.
- Reusable mocks belong near the owning domain interfaces.

## Backend Guardrails

- Do not leak provider-specific language into `domain/` or `useCases/`.
- Do not duplicate transport DTOs when the canonical use-case DTO already fits.
- Do not scatter route registration or business rules across handlers.
- If a new backend convention is needed, ask the user instead of inventing it.
