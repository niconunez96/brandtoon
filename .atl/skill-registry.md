# Skill Registry

Project: brandtoon
Generated: 2026-04-08
Mode: base SDD conventions + Engram persistence

## Resolution Rules

- Project-level skills override user-level skills with the same name.
- `sdd-*`, `_shared`, and `skill-registry` skills are excluded from this registry.
- Root integration work follows `AGENTS.md` first, then scoped guides under `front/` and `back/`.

## Project Conventions

| Type | Path | Purpose |
|---|---|---|
| gateway | `AGENTS.md` | Root routing for frontend, backend, and integration work |
| scope | `front/AGENTS.md` | Frontend baseline, state boundaries, required frontend skills |
| scope | `back/AGENTS.md` | Backend DDD + CQRS + Huma + TDD governance |

## Project Skills

| Skill | Level | Path | Trigger / Use When |
|---|---|---|---|
| justfile-expert | project | `.agents/skills/justfile-expert/SKILL.md` | justfile, recipes, task automation, project commands |
| brandtoon-frontend-baseline | project | `front/.agents/skills/brandtoon-frontend-baseline/SKILL.md` | any feature work in `front/`; frontend architecture and state boundaries |
| frontend-design | project | `front/.agents/skills/frontend-design/SKILL.md` | explicit visual redesign or polished frontend/artifact creation |
| tailwind-design-system | project | `front/.agents/skills/tailwind-design-system/SKILL.md` | Tailwind design systems, tokens, standardized UI patterns |
| vercel-react-best-practices | project | `front/.agents/skills/vercel-react-best-practices/SKILL.md` | React code, performance work, data fetching, bundle optimization |
| web-design-guidelines | project | `front/.agents/skills/web-design-guidelines/SKILL.md` | UI review, accessibility audit, UX/design best-practice checks |
| go-ddd-aggregate | project | `back/.agents/skills/go-ddd-aggregate/SKILL.md` | editing `back/bounded_contexts/**/domain` |
| go-usecase-cqrs-functions | project | `back/.agents/skills/go-usecase-cqrs-functions/SKILL.md` | editing `back/bounded_contexts/**/useCases` |
| go-infra-http-repo-adapters | project | `back/.agents/skills/go-infra-http-repo-adapters/SKILL.md` | editing `back/bounded_contexts/**/infra` |
| go-huma-api-contracts | project | `back/.agents/skills/go-huma-api-contracts/SKILL.md` | Huma API contracts, validation, `back/main.go`, `infra/http` |
| go-shared-di-container | project | `back/.agents/skills/go-shared-di-container/SKILL.md` | `bounded_contexts/shared/diContainer.go`, composition root wiring |
| go-testing-tdd-backend | project | `back/.agents/skills/go-testing-tdd-backend/SKILL.md` | backend `useCases` / `infra/http` changes requiring RED-GREEN-REFACTOR |

## User Skills Available

| Skill | Level | Path | Trigger / Use When |
|---|---|---|---|
| branch-pr | user | `~/.config/opencode/skills/branch-pr/SKILL.md` | creating or preparing a pull request |
| issue-creation | user | `~/.config/opencode/skills/issue-creation/SKILL.md` | creating a GitHub issue |
| go-testing | user | `~/.config/opencode/skills/go-testing/SKILL.md` | Go tests and Bubbletea TUI testing |
| judgment-day | user | `~/.config/opencode/skills/judgment-day/SKILL.md` | adversarial dual review |
| skill-creator | user | `~/.config/opencode/skills/skill-creator/SKILL.md` | creating new AI skills |
| workito-ai-git-workflow | user | `~/.config/opencode/skills/workito-ai-git-workflow/SKILL.md` | branch/worktree/commit/PR workflow orchestration |

## Default Stack Skill Resolution

### Frontend (`front/`)

1. `brandtoon-frontend-baseline`
2. `vercel-react-best-practices`
3. `web-design-guidelines`
4. `frontend-design` only when explicitly requested
5. `tailwind-design-system` when design system or Tailwind standardization work appears

### Backend (`back/`)

Load by touched layer:

- `go-ddd-aggregate` -> domain modeling
- `go-usecase-cqrs-functions` -> use case logic
- `go-infra-http-repo-adapters` -> infra/http or repo adapters
- `go-huma-api-contracts` -> API contracts and Huma registration
- `go-shared-di-container` -> dependency wiring and container work
- `go-testing-tdd-backend` -> any `useCases` or `infra/http` change

### Integration (root)

1. Follow `AGENTS.md` contract-first order.
2. Backend contract is source of truth.
3. Then load frontend/backend scoped skills based on touched files.
