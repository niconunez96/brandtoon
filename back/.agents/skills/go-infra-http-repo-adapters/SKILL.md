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

## Canonical Policy First

- Backend policy lives in `back/AGENTS.md`.
- This skill adds infra-layer adapter, routing, and repository guidance.

## Objective

- Keep third-party integrations isolated in infra.
- Standardize route registration and handler boundaries.
- Keep adapter naming explicit and searchable.

## Hard Constraints

- Infra includes `repo/` and `http/` directories.
- Aggregate package naming is mandatory:
  - `infra/repo/` files use `package {aggregate}repo`
  - `infra/http/` files use `package {aggregate}http`
- Repository implementations use `XXXPostgresRepo` naming.
- HTTP routes are grouped by prefix in `infra/http/routes.go`.
- Handlers are function-based and delegate business logic to use cases.
- Aggregate-owned endpoints MUST be registered in the same aggregate's `infra/http` layer as the use case they invoke; do not host an endpoint in a foreign aggregate just because the URL path mentions that concept.
- Transport payload DTOs must come from the owning aggregate's `useCases` layer; `infra/http` should use them directly when the shape matches and must not redefine the same payload shape as a parallel DTO.
- Infra adapters depend inward on domain contracts/use cases.
- Dependency wiring standards are defined in `go-shared-di-container`.
- Provider/vendor-specific naming is ALLOWED in infra adapters and route paths when mapping concrete integrations (for example: `GoogleOAuthClient`, `/auth/google/callback`).

## Forbidden Patterns

- Business rules implemented in HTTP handlers.
- Endpoint registration in one aggregate that delegates aggregate-owned behavior to another aggregate's use case.
- Defining duplicate response/request DTO structs in `infra/http` when the same payload already exists in `useCases`.
- Direct domain mutations from SQL result mapping without aggregate methods.
- Multiple route roots spread across endpoint files.
- Ambiguous adapter names (`RepoImpl`, `StorageAdapter`).
- Generic infra package names inside aggregate layers (`package repo`, `package http`).
- Domain layer importing HTTP, SQL, orm, or third-party clients.
- Initializing external clients directly inside handlers or useCases.
- Re-instantiating the same third-party dependency repeatedly instead of reusing container-managed singleton instances.

## Required Checklist

- Confirm all integrations stay under infra.
- Validate route prefix grouping in `routes.go`.
- Keep one handler function per endpoint concern.
- Ensure handlers map DTOs and call use cases only.
- Reuse use-case DTOs for transport payloads; only add aliasing/mapping when transport-specific concerns truly differ, and prefer direct references over alias files.
- Ensure repository structs satisfy domain interfaces explicitly.
- Ensure infra package names follow `{aggregate}repo` and `{aggregate}http`.
- Ensure dependency wiring changes follow `go-shared-di-container`.
- Confirm provider-specific terms stay contained to `infra/` and do not leak into `domain/` or `useCases/` naming.
- After editing backend infra code, run `just back-format` and then `just back-check`; iterate until the canonical backend check passes.

## Output Standard

- Adapter boundaries are obvious from folder and naming conventions.
- Dependency direction remains `main -> routes -> handlers -> useCases -> domain`.
- New integrations are replaceable without domain changes.

See `REFERENCE.md` for concrete layout and examples.
