---
name: go-testing-tdd-backend
description: |
  Enforces backend TDD workflow in Brandtoon for changes in useCases and HTTP handlers.
  Requires RED-GREEN-REFACTOR, test completion gates, and reusable mocks near domain interfaces.
user-invocable: false
allowed-tools: Bash, Grep, Glob, Read, Write, Edit
---

# Go Backend Testing (TDD)

Use this skill when changing `back/bounded_contexts/**/useCases` or `back/bounded_contexts/**/infra/http`.

## Objective

- Enforce fast feedback loop with strict TDD.
- Keep tests focused on use cases and HTTP handlers.
- Reuse mocks across tests by colocating them near domain interfaces.

## Mandatory TDD Flow

1. **RED**: add/update a failing test first.
2. **GREEN**: implement the minimum code required for passing tests.
3. **REFACTOR**: improve code while keeping tests green.

Do NOT finish a feature/fix/refactor until all relevant tests pass.

## Test Scope (Required)

- Required tests for touched code in:
  - `bounded_contexts/{context}/{aggregate}/useCases`
  - `bounded_contexts/{context}/{aggregate}/infra/http`
- If tests are missing in touched area, add baseline tests before marking done.

## Test Organization

- Use case tests:
  - `bounded_contexts/{context}/{aggregate}/useCases/test/*_test.go`
- HTTP handler tests:
  - `bounded_contexts/{context}/{aggregate}/infra/http/test/*_test.go`

Keep one clear intent per test case. Prefer outcome assertions over internal behavior assertions.

## Mock Reuse Policy (Hard Rule)

- Put mocks near the domain interface they implement:
  - Aggregate-specific mocks: `bounded_contexts/{context}/{aggregate}/domain/mocks/`
  - Shared mocks: `bounded_contexts/shared/domain/mocks/`
- Reuse existing domain mocks across use case and handler tests.
- Do not create repository/service mocks inline in individual test files when reusable placement exists.

## Handler Testing Rules

- Test Huma handlers through registered operations using `humatest` (preferred).
- `httptest` is acceptable when `humatest` does not fit a specific case.
- Keep business logic assertions at use case/domain level.
- Handler tests focus on transport mapping, status codes, and payload contracts.

## Huma Testing Reference

- Follow Huma tutorial approach: extract route registration to reusable function and mount it in tests.
- Source: https://huma.rocks/tutorial/writing-tests/

## Completion Gate

A backend change is complete only if:

- Tests exist for touched useCases/handlers.
- New/updated tests pass.
- No relevant failing tests remain.
- Missing reusable mocks were created in the correct `domain/mocks` location.

See `REFERENCE.md` for templates and practical patterns.
