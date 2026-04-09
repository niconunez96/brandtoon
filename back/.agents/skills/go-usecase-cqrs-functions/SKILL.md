---
name: go-usecase-cqrs-functions
description: |
  Enforces function-based use cases with CQRS semantics in backend Go code.
  One file per use case, explicit command/query separation, and dependency injection
  through function arguments and domain interfaces.
user-invocable: false
allowed-tools: Bash, Grep, Glob, Read, Write, Edit
---

# Go Use Cases (CQRS Functions)

Use this skill when editing `back/bounded_contexts/**/useCases`.

## Objective

- Keep application logic explicit, composable, and testable.
- Enforce command/query separation by behavior.
- Keep dependency injection visible through function signatures.

## Hard Constraints

- One file per use case (`authenticate_user.go`, `reset_password.go`, ...).
- Use cases are functions, not service structs.
- Each use case receives:
  - command/query input DTO
  - required domain interfaces as parameters
- Commands perform writes or side effects.
- Queries return read models without side effects.

## Forbidden Patterns

- Multi-use-case files mixing unrelated flows.
- Stateful use-case structs storing mutable dependencies.
- Commands that silently behave as queries (or inverse).
- Use cases importing transport concerns (HTTP DTOs, router context).
- Use cases coupling directly to infra implementations.

## Required Checklist

- Name file by business action in snake_case.
- Define input type per use case (command or query).
- Pass only required dependencies in function signature.
- Return domain/application errors explicitly.
- Confirm command/query intent by side-effect behavior.

## Output Standard

- Use case functions are focused and single-purpose.
- Signatures reveal contract and dependencies immediately.
- Handlers can orchestrate without embedding business rules.

See `REFERENCE.md` for canonical templates.
