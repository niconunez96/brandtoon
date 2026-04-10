---
name: go-ddd-aggregate
description: |
  Backend DDD aggregate modeling for Brandtoon Go codebase. Enforces aggregate root
  boundaries, value objects, encapsulated entities, and pure domain contracts.
user-invocable: false
allowed-tools: Bash, Grep, Glob, Read, Write, Edit
---

# Go DDD Aggregate

Use this skill when changing domain models in `back/bounded_contexts/**/domain`.

## Objective

- Keep domain behavior inside aggregate root boundaries.
- Protect invariants through explicit methods and value objects.
- Keep external coupling behind domain interfaces only.

## Hard Constraints

- Aggregate root is the ONLY entrypoint for domain interactions.
- Internal entities are not exported for direct external mutation.
- Value objects are defined via explicit custom types and enum-like constants.
- Domain interfaces (ports) are declared in domain layer and are third-party free.
- Domain must compile without infra/framework imports.
- Domain names must be provider-agnostic and business-language first.
- Domain identifiers MUST NOT include provider/vendor terms (for example: `Google`, `GitHub`, `Auth0`).

## Forbidden Patterns

- Anemic models with setters everywhere.
- Domain services returning infra-specific types.
- External code mutating child entities directly.
- Implicit string literals instead of typed value objects.
- Interfaces defined in infra for domain dependencies.
- Provider-branded domain contracts or entities (`GoogleIdentityProvider`, `GithubUser`, etc.).
- Provider-specific domain fields that leak third-party semantics (`googleSubject`, `githubLogin`).

## Required Checklist

- Identify aggregate root and invariants before coding.
- Confirm all mutations pass through aggregate root methods.
- Replace primitive obsession with typed value objects where needed.
- Define/adjust domain interfaces for external side effects.
- Verify no third-party dependency leaks into domain signatures.
- Verify domain names survive provider swap without renaming (`Google -> GitHub` test).

## Output Standard

- Domain code remains deterministic and testable in isolation.
- Naming expresses ubiquitous language of each bounded context.
- Contracts are minimal and intention-revealing.

See `REFERENCE.md` for templates and examples.
