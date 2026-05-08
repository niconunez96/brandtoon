# Agent Policy

Repo-wide policy for agents in `brandtoon`.

## Ownership

- Root `AGENTS.md` owns routing between `front/`, `back/`, and integration work.
- `front/AGENTS.md` owns frontend-only policy.
- `back/AGENTS.md` owns backend-only policy.
- Skills own operational workflow, checklists, and step-by-step implementation guidance.

## Precedence

1. Direct user instruction
2. Nearest scoped `AGENTS.md`
3. Root `AGENTS.md`
4. These standards

## Required Behavior

- Keep changes minimal, explicit, and correctly encapsulated.
- Prefer references over duplicated guidance.
- Preserve documented architecture boundaries.
- Do not invent new policy when the repo is silent.

## Ask the User When

- The convention is unclear or missing.
- Frontend/backend ownership is ambiguous.
- A change would establish a new project rule.
- Two valid interpretations remain after reading the repo guidance.

When blocked by ambiguity, ask the user instead of improvising.
