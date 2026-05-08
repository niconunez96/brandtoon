# Skill Authoring Standards

Skills are the operational layer for this repository.

## What Skills Should Contain

- Execution workflow
- Validation and testing loops
- Detailed checklists
- Examples, templates, and references

## What Skills Should Not Own

- Scope routing between `front/`, `back/`, and root
- Repo-wide precedence rules
- Duplicated project policy already stated in canonical standards or scoped `AGENTS.md`

## Authoring Rules

- Keep project policy in `AGENTS.md`; keep operational detail in skills.
- Reference canonical docs instead of restating them when possible.
- Scope each skill clearly by path, task, or trigger.
- Update a skill when workflow guidance changes.
- If a workflow depends on an unstated convention, instruct the agent to ask the user rather than improvise.
