---
name: go-huma-api-contracts
description: |
  Huma framework conventions for Brandtoon backend HTTP APIs. Use for service configuration,
  operation registration, request/response contracts, validation tags, and Huma-specific testing setup.
user-invocable: false
allowed-tools: Bash, Grep, Glob, Read, Write, Edit
---

# Go Huma API Contracts

Use this skill when changing `back/main.go` or `back/bounded_contexts/**/infra/http`.

## Objective

- Keep HTTP APIs aligned with Huma contracts and validation model.
- Keep route/operation registration testable.
- Preserve DDD flow: handlers/use cases/domain boundaries.

## Required Patterns

1. Use Huma operation registration patterns (`huma.Register` / `huma.Get` etc.) consistently.
2. Prefer explicit request/response structs with schema tags for docs + validation.
3. Default approach for input schema validation is Huma field tags in `infra/http` contracts (e.g. `minLength`, `maxLength`, `pattern`) instead of duplicating the same schema checks in use cases.
4. Keep route registration isolated and callable from tests.
5. Keep service configuration externalized (CLI/env options) when needed.
6. Use conditional field constraints with `dependentRequired` when business contracts need field dependencies.
7. For dependency initialization/composition wiring, follow `go-shared-di-container`.

## DDD Boundary Rules

- Huma-specific types stay in transport (`infra/http`) and composition root.
- Do not leak Huma request/response structs into domain.
- Map Huma input/output to use case command/query DTOs explicitly.
- Keep transport-level validation failures as HTTP validation concerns; use cases should focus on business behavior and invariants.

## Testing Alignment

- HTTP endpoint tests should use Huma's `humatest` utilities where practical.
- Keep route registration in a reusable function so tests can mount API and call endpoints.
- Transport tests assert status code, payload, and validation behavior.

## Forbidden

- Coupling domain logic to Huma internals.
- Inlining all registration inside `main` with no reusable route wiring function.
- Skipping Huma validation/schema tags when endpoint contracts need them.

See `REFERENCE.md` for concrete patterns based on Huma docs.
