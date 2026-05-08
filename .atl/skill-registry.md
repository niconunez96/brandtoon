# Skill Registry

Project: brandtoon
Generated: 2026-05-08

## Canonical Docs

- Root gateway: `AGENTS.md`
- Repo-wide agent policy: `.agents/standards/agent-policy.md`
- Skill authoring policy: `.agents/standards/skill-authoring.md`
- Frontend scope: `front/AGENTS.md`
- Backend scope: `back/AGENTS.md`

## Resolution Rules

- Project-level skills override user-level skills with the same name.
- `sdd-*`, `_shared`, and `skill-registry` skills are excluded from runtime matching here.
- Keep project policy in `AGENTS.md` and standards docs; keep operational detail in skills.

## Project Skills

| Skill | Path | Use when |
|---|---|---|
| justfile-expert | `.agents/skills/justfile-expert/SKILL.md` | Justfile, recipes, task automation |
| brandtoon-frontend-baseline | `front/.agents/skills/brandtoon-frontend-baseline/SKILL.md` | Any work in `front/` needing frontend workflow and boundaries |
| frontend-design | `front/.agents/skills/frontend-design/SKILL.md` | Explicit visual redesign |
| tailwind-design-system | `front/.agents/skills/tailwind-design-system/SKILL.md` | Design tokens or component standardization |
| vercel-react-best-practices | `front/.agents/skills/vercel-react-best-practices/SKILL.md` | React performance and implementation patterns |
| web-design-guidelines | `front/.agents/skills/web-design-guidelines/SKILL.md` | UI review, accessibility, UX audits |
| go-ddd-aggregate | `back/.agents/skills/go-ddd-aggregate/SKILL.md` | `back/**/domain` modeling |
| go-usecase-cqrs-functions | `back/.agents/skills/go-usecase-cqrs-functions/SKILL.md` | `back/**/useCases` work |
| go-infra-http-repo-adapters | `back/.agents/skills/go-infra-http-repo-adapters/SKILL.md` | `back/**/infra` work |
| go-huma-api-contracts | `back/.agents/skills/go-huma-api-contracts/SKILL.md` | Huma contracts, validation, route wiring |
| go-shared-di-container | `back/.agents/skills/go-shared-di-container/SKILL.md` | Shared dependency wiring and `main.go` composition |
| go-testing-tdd-backend | `back/.agents/skills/go-testing-tdd-backend/SKILL.md` | Backend TDD for `useCases` and `infra/http` |

## User Skills Available

| Skill | Path | Use when |
|---|---|---|
| workito-ai-git-workflow | `~/.config/opencode/skills/workito-ai-git-workflow/SKILL.md` | Branch/worktree/commit/PR workflow |
| branch-pr | `~/.config/opencode/skills/branch-pr/SKILL.md` | Creating or preparing a PR |
| issue-creation | `~/.config/opencode/skills/issue-creation/SKILL.md` | Creating GitHub issues |
| skill-creator | `~/.config/opencode/skills/skill-creator/SKILL.md` | Creating new skills |
| go-testing | `~/.config/opencode/skills/go-testing/SKILL.md` | Go testing patterns |
| judgment-day | `~/.config/opencode/skills/judgment-day/SKILL.md` | Dual review workflow |
| find-skills | `~/.agents/skills/find-skills/SKILL.md` | Discovering/installing skills |

## Scoped Skill Defaults

### Frontend

1. `brandtoon-frontend-baseline`
2. `vercel-react-best-practices`
3. `web-design-guidelines`
4. `frontend-design` only when explicitly requested
5. `tailwind-design-system` when design-system standardization is part of the task

### Backend

- `go-ddd-aggregate` -> domain
- `go-usecase-cqrs-functions` -> use cases
- `go-infra-http-repo-adapters` -> infra
- `go-huma-api-contracts` -> Huma transport contracts
- `go-shared-di-container` -> DI wiring
- `go-testing-tdd-backend` -> `useCases` or `infra/http`

### Integration

1. Follow `AGENTS.md` and the canonical standards first.
2. Backend contract remains the source of truth.
3. Then load frontend/backend scoped skills by touched files.
