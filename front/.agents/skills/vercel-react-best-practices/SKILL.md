---
name: vercel-react-best-practices
description: Brandtoon-compatible React and Vite performance guidance distilled from the broader Vercel React guidance. Use for React component work, client-side async flows, bundle trimming, and render performance in front/.
license: MIT
metadata:
  author: vercel
  version: "1.0.0"
---

# Brandtoon React + Vite Best Practices

Use this skill for `front/` React work that needs BETTER performance or lower incidental complexity.

## Canonical Policy First

- Frontend policy lives in `front/AGENTS.md`.
- Frontend workflow/state boundaries live in `front/.agents/skills/brandtoon-frontend-baseline/SKILL.md`.
- This skill adds **React/Vite implementation guidance only**.

## Brandtoon Scope

Apply these rules to:
- React components and hooks
- Vite bundle/runtime decisions
- Client-side async flows
- Local render and interaction performance

Do **not** pull in Next.js-, RSC-, Server Actions-, or SWR-specific guidance unless the user explicitly changes the stack.

## Keep These Rules

### 1. Remove avoidable async waterfalls

- Check cheap synchronous guards before awaiting.
- Start independent promises early and join with `Promise.all()`.
- Defer awaits into the branch that actually needs the data.

### 2. Trim the client bundle

- Lazy-load heavy UI with `React.lazy()` or dynamic `import()`.
- Load optional modules only when the feature is activated.
- Preload likely-next heavy modules on hover/focus when it materially improves UX.
- Prefer direct imports or proven package-level optimizations when barrel imports hurt Vite dev/build performance.

### 3. Reduce unnecessary re-renders

- Derive state during render when it can be computed from current props/state.
- Put interaction-triggered side effects in event handlers, not effects.
- Use functional `setState` updates when next state depends on previous state.
- Do not define components inside components.
- Use `useMemo` only for genuinely expensive derived work, not simple primitive expressions.
- Split unrelated effects/memos so dependency changes do not recompute everything.

### 4. Keep input and scrolling responsive

- Use `useDeferredValue` or `startTransition` for non-urgent expensive updates.
- Use passive listeners for scroll/touch handlers that do not call `preventDefault()`.
- Store transient non-visual values in refs when re-rendering is unnecessary.

### 5. Prefer immutable and low-overhead data transforms

- Use `toSorted()` or copied arrays instead of mutating props/state with `sort()`.
- Use `Map`/`Set` for repeated lookups.
- Combine repeated iterations in hot paths when the code stays readable.
- Return early when the outcome is already known.

### 6. Treat browser storage as an I/O boundary

- Minimize and version `localStorage` payloads.
- Wrap storage access in `try/catch`.
- Do not store server payloads wholesale when the UI only needs a few fields.

## Brandtoon-Specific Notes

- Prefer TanStack Query for backend-backed async state; do not import SWR guidance from upstream Vercel docs.
- Prefer plain React/Vite primitives over framework-specific abstractions.
- If a recommendation conflicts with `front/AGENTS.md` or `front/docs/design-foundations.md`, the local Brandtoon docs win.

## Reference Material

Use the local `rules/` files as examples only after confirming the rule still fits the Vite + React stack.
