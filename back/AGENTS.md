# Backend AI Governance (Scope: `back/` only)

This file governs ALL AI-driven work inside `back/`.

## Mission

- Preserve a strict DDD + CQRS-by-use-case structure.
- Keep domain pure and stable.
- Keep infrastructure replaceable.
- Enforce functional use-case and handler style.

## Mandatory Backend Structure

- Base path: `bounded_contexts/{context}/{aggregate}/`
- Layers per aggregate: `domain/`, `infra/`, `useCases/`
- Shared modules: `bounded_contexts/shared/domain/`, `bounded_contexts/shared/infra/`
- Shared dependency wiring: `bounded_contexts/shared/diContainer.go` (DI container)

Expected shape:

```text
bounded_contexts/
  shared/
    domain/
    infra/
  identity/
    user/
      domain/
      useCases/
      infra/
```

## Dependency Flow (Hard Rule)

- `main.go -> routes.go -> handlers -> useCases -> domain`
- `main.go` owns dependency initialization and injection.
- `main.go` must resolve dependencies from shared `DIContainer` getters and pass them down.
- Domain NEVER imports infra or transport.

## Auto-Load Skills (BEFORE coding)

When work touches `back/`, load skills by context before writing code:

- Aggregate/domain modeling -> `go-ddd-aggregate`
- Application behavior/CQRS use cases -> `go-usecase-cqrs-functions`
- HTTP, repository adapters, integrations -> `go-infra-http-repo-adapters`
- Any Huma API wiring, operation contracts, service config, validation, or schema tags -> `go-huma-api-contracts`
- Any dependency initialization/composition root wiring (`diContainer.go`, `main.go`) -> `go-shared-di-container`
- Any feature/fix/refactor touching `useCases` or `infra/http` -> `go-testing-tdd-backend`

If work spans multiple areas, load multiple skills.

## Hard Constraints

- One use case per file (`authenticate_user.go`, `reset_password.go`, ...).
- Use cases are functions (not service structs) that receive:
  - command/query DTO
  - injected domain interfaces as function args
- Commands perform side effects; queries perform read operations.
- Handlers are functions (not stateful classes).
- HTTP routes are grouped by prefix in `infra/http/routes.go`.
- Huma is the standard HTTP framework for `back/`.
- Repository adapters live in `infra/repo/`.
- Postgres repositories follow `XXXPostgresRepo` naming.
- Third-party dependency initialization must live in shared `DIContainer` (not in handlers/useCases/domain).
- `DIContainer` exposes `GetX()` methods that lazily initialize once (singleton-style) and reuse instances.
- Domain contracts define external ports (`EmailSender`, `XXRepository`) with NO third-party dependencies in signatures.
- Aggregate root is the only domain interaction entrypoint.
- Internal entities stay encapsulated inside aggregate boundaries.
- Value objects are modeled as explicit enum/type definitions.
- TDD is mandatory for backend changes in `useCases` and `infra/http`:
  - RED: write/adjust a failing test first
  - GREEN: implement minimal code to pass
  - REFACTOR: improve design while tests remain green
- A backend feature/fix/refactor is NOT complete until all relevant tests pass.
- If tests are missing for touched `useCases` or `infra/http` handlers, add them.
- Tests are required only for:
  - `bounded_contexts/{context}/{aggregate}/useCases`
  - `bounded_contexts/{context}/{aggregate}/infra/http`
- Test file placement:
  - `bounded_contexts/{context}/{aggregate}/useCases/test/*_test.go`
  - `bounded_contexts/{context}/{aggregate}/infra/http/test/*_test.go`
- Mocks must be reusable and placed closest to their domain interface:
  - Aggregate-specific: `bounded_contexts/{context}/{aggregate}/domain/mocks/`
  - Shared cross-context: `bounded_contexts/shared/domain/mocks/`
- Never define repository/service mocks inline in test files when they can live in `domain/mocks`.

## Forbidden Patterns

- Domain importing framework/driver packages (`gin`, `fiber`, `gorm`, `sqlx`, SMTP clients, etc.).
- Use case structs with hidden mutable state.
- Fat handlers containing business rules.
- Route definitions scattered across handler files.
- Generic repository names that hide implementation (`UserRepoImpl`).
- Cross-aggregate mutation bypassing aggregate roots.
- Marking a change as done with failing tests.
- Skipping test creation for modified `useCases` or `infra/http` handlers.
- Duplicating one-off test-local mocks instead of reusing `domain/mocks`.
- Implementing new HTTP API endpoints without following Huma operation/contract patterns.
- Initializing third-party clients directly inside handlers, useCases, or domain.
- Bypassing `DIContainer` and constructing infra dependencies ad-hoc in `main.go`.

## Execution Checklist (Use on every backend task)

- Confirm paths are only under `back/`.
- Load required skill(s) before editing.
- Validate dependency direction stays inward to domain.
- Validate names and file placement follow conventions.
- Ensure commands/queries are separated clearly.
- Ensure domain contracts stay vendor-neutral.

## Validation Commands

- Canonical backend validation entrypoint: `just back-check`
- Non-mutating backend validation steps:
  - `just back-format-check`
  - `just back-lint`
  - `just back-test`
- CI must stay non-mutating: do NOT use `just back-format` in validation workflows.
