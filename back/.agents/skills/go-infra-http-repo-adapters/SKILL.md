---
name: go-infra-http-repo-adapters
description: |
  Enforces backend infrastructure adapter conventions for HTTP transport and repository
  implementations in Go. Covers route grouping, handler functions, and Postgres adapter naming.
user-invocable: false
allowed-tools: Bash, Grep, Glob, Read, Write, Edit
---

# Go Infra (HTTP + Repo Adapters)

Use this skill when editing `back/bounded_contexts/**/infra`.

## Objective

- Keep third-party integrations isolated in infra.
- Standardize route registration and handler boundaries.
- Keep adapter naming explicit and searchable.

## Hard Constraints

- Infra includes `repo/` and `http/` directories.
- Repository implementations use `XXXPostgresRepo` naming.
- HTTP routes are grouped by prefix in `infra/http/routes.go`.
- Handlers are function-based and delegate business logic to use cases.
- Infra adapters depend inward on domain contracts/use cases.
- Dependency wiring standards are defined in `go-shared-di-container`.
- Provider/vendor-specific naming is ALLOWED in infra adapters and route paths when mapping concrete integrations (for example: `GoogleOAuthClient`, `/auth/google/callback`).

## Forbidden Patterns

- Business rules implemented in HTTP handlers.
- Direct domain mutations from SQL result mapping without aggregate methods.
- Multiple route roots spread across endpoint files.
- Ambiguous adapter names (`RepoImpl`, `StorageAdapter`).
- Domain layer importing HTTP, SQL, orm, or third-party clients.
- Initializing external clients directly inside handlers or useCases.
- Re-instantiating the same third-party dependency repeatedly instead of reusing container-managed singleton instances.

## Required Checklist

- Confirm all integrations stay under infra.
- Validate route prefix grouping in `routes.go`.
- Keep one handler function per endpoint concern.
- Ensure handlers map DTOs and call use cases only.
- Ensure repository structs satisfy domain interfaces explicitly.
- Ensure dependency wiring changes follow `go-shared-di-container`.
- Confirm provider-specific terms stay contained to `infra/` and do not leak into `domain/` or `useCases/` naming.

## Output Standard

- Adapter boundaries are obvious from folder and naming conventions.
- Dependency direction remains `main -> routes -> handlers -> useCases -> domain`.
- New integrations are replaceable without domain changes.

See `REFERENCE.md` for concrete layout and examples.
