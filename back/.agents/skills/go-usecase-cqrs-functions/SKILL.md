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
- Package name in `useCases/` is `package {aggregate}usecases`.
- Use cases are functions, not service structs.
- Each use case receives:
  - command/query input DTO
  - required domain interfaces as parameters
- Commands perform writes or side effects.
- Queries return read models without side effects.
- Use case file/function/DTO names must express business capability, not provider technology.
- Avoid transport schema validation in use cases (required/length/pattern checks for HTTP payload shape); prefer Huma validation tags in `infra/http` contracts.

## Forbidden Patterns

- Multi-use-case files mixing unrelated flows.
- Generic package names in aggregate use-case layers (`package useCases`, `package usecases`).
- Stateful use-case structs storing mutable dependencies.
- Commands that silently behave as queries (or inverse).
- Use cases importing transport concerns (HTTP DTOs, router context).
- Use cases coupling directly to infra implementations.
- Duplicating schema validation rules both in Huma contracts and in use cases.
- Provider-branded use case names (`AuthenticateGoogleCallback`, `GetGithubAuthURL`, etc.).
- Provider-specific command/query DTO field names in `useCases/` (`googleSubject`, `githubCode`).

## Required Checklist

- Name file by business action in snake_case.
- Ensure package name matches `{aggregate}usecases`.
- Define input type per use case (command or query).
- Pass only required dependencies in function signature.
- Return domain/application errors explicitly.
- Confirm command/query intent by side-effect behavior.
- Apply provider-swap naming check: would `Google -> GitHub` require renaming in `useCases/`? If yes, rename to domain language.

## Output Standard

- Use case functions are focused and single-purpose.
- Signatures reveal contract and dependencies immediately.
- Handlers can orchestrate without embedding business rules.

See `REFERENCE.md` for canonical templates.
