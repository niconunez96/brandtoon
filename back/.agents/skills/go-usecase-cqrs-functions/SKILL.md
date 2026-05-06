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
- Use cases MUST return use-case DTOs, not domain entities/aggregates.
- Use-case DTOs are the canonical payload types for the aggregate; downstream `infra/http` code should reference them directly when the shape matches and must not introduce parallel DTO structs for the same payload.
- Use case file/function/DTO names must express business capability, not provider technology.
- Use cases MUST live in the aggregate that owns the mutated/read model. Do not place cross-aggregate orchestration use cases in foreign aggregates.
- If a use case owns aggregate behavior, its HTTP endpoint/handler must be hosted by the same aggregate's `infra/http` layer.
- DTO naming convention is mandatory: file `{name}_dto.go`, struct `{Name}DTO`.
- Each DTO file must expose `serialize` (single domain object -> DTO) and `serializeList` (list of domain objects -> list of DTOs).
- Avoid alias files for DTOs unless there is a concrete transport reason; direct references are preferred for clarity.
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
- DTO transformations are explicit and follow `{name}_dto.go` + `{Name}DTO` + `serialize` / `serializeList` convention.

See `REFERENCE.md` for canonical templates.
