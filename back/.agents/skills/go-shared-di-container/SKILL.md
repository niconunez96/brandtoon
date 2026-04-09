---
name: go-shared-di-container
description: |
  Shared dependency container conventions for Brandtoon backend. Use when editing
  `bounded_contexts/shared/diContainer.go` or wiring dependencies from `main.go`.
user-invocable: false
allowed-tools: Bash, Grep, Glob, Read, Write, Edit
---

# Go Shared DI Container

Use this skill when changing dependency initialization or composition root wiring.

## Scope

- `back/bounded_contexts/shared/diContainer.go`
- `back/main.go`
- Dependency composition code for route/handler setup

## Hard Constraints

- Third-party dependency initialization must live in `diContainer.go`.
- Container exposes `GetX()` methods with lazy singleton-style initialization.
- `main.go` resolves dependencies via `NewDIContainer()` + getter calls.
- Handlers/useCases/domain must not initialize third-party clients directly.
- Getter return types should prefer domain interfaces when possible.

## Canonical Pattern

1. Keep private fields inside `DIContainer` for cached instances.
2. In `GetX()`, initialize only when nil.
3. Reuse shared clients across multiple getters (e.g. DB, HTTP client).
4. Compose route dependencies in `main.go` by calling getters.

## Forbidden Patterns

- Ad-hoc dependency construction in `main.go` bypassing container getters.
- Re-initializing heavy resources per request or per handler.
- Leaking concrete infra types into domain layer signatures.
- Circular dependency creation inside container getters.

## Completion Checklist

- All third-party clients are created in `diContainer.go`.
- Getters are idempotent and reuse the same instance.
- `main.go` only composes through container getters.
- Dependency direction stays `main -> routes -> handlers -> useCases -> domain`.

See `REFERENCE.md` for examples.
