# React Best Practices

This directory is an upstream Vercel React performance reference that Brandtoon keeps as supporting material only.

## Brandtoon Usage

- Canonical repo behavior lives in `../../AGENTS.md` and `../../../.agents/standards/agent-policy.md`.
- Canonical operational guidance for this skill lives in `SKILL.md`.
- The broad upstream `rules/` content is **not** authoritative for Brandtoon when it conflicts with the Vite + React + TanStack Query stack.
- Do not adopt Next.js-, SWR-, Server Actions-, or React-Compiler-specific patterns here unless the user explicitly changes the project stack.

## Structure

- `rules/` - Individual rule files (one per rule)
  - `_sections.md` - Section metadata (titles, impacts, descriptions)
  - `_template.md` - Template for creating new rules
  - `area-description.md` - Individual rule files
- `src/` - Build scripts and utilities
- `metadata.json` - Document metadata (version, organization, abstract)
- __`AGENTS.md`__ - Compiled output (generated)
- __`test-cases.json`__ - Test cases for LLM evaluation (generated)

## Brandtoon Maintenance Rule

- Do **not** treat this directory as an active build/generated workflow for routine Brandtoon work.
- `AGENTS.md`, `test-cases.json`, and the `src/` scripts belong to the upstream reference toolchain, not the default Brandtoon operating model.
- If this upstream mirror ever needs maintenance, first verify the change is explicitly required and then follow the local `SKILL.md` plus canonical standards.

## Rule File Structure

Each rule file should follow this structure:

```markdown
---
title: Rule Title Here
impact: MEDIUM
impactDescription: Optional description
tags: tag1, tag2, tag3
---

## Rule Title Here

Brief explanation of the rule and why it matters.

**Incorrect (description of what's wrong):**

```typescript
// Bad code example
```

**Correct (description of what's right):**

```typescript
// Good code example
```

Optional explanatory text after examples.

Reference: [Link](https://example.com)

## File Naming Convention

- Files starting with `_` are special (excluded from build)
- Rule files: `area-description.md` (e.g., `async-parallel.md`)
- Section is automatically inferred from filename prefix
- Rules are sorted alphabetically by title within each section
- IDs (e.g., 1.1, 1.2) are auto-generated during build

## Impact Levels

- `CRITICAL` - Highest priority, major performance gains
- `HIGH` - Significant performance improvements
- `MEDIUM-HIGH` - Moderate-high gains
- `MEDIUM` - Moderate performance improvements
- `LOW-MEDIUM` - Low-medium gains
- `LOW` - Incremental improvements

## Upstream Tooling Reference

These scripts describe the upstream source project only:

- `pnpm build` - Compile rules into AGENTS.md
- `pnpm validate` - Validate all rule files
- `pnpm extract-tests` - Extract test cases for LLM evaluation
- `pnpm dev` - Build and validate

## Contributing

For Brandtoon, prefer updating `SKILL.md`, scoped `AGENTS.md`, or canonical standards when local guidance changes.
Only touch this upstream reference structure when the task is specifically about maintaining the imported Vercel material.

## Acknowledgments

Originally created by [@shuding](https://x.com/shuding) at [Vercel](https://vercel.com).
