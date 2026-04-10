# Brandtoon Design Foundations

This document is the reusable source of truth for the first frontend design-system layer.

## 1. Theme tokens

Defined in `src/styles.css`.

### Colors

- `bg-page` — application canvas
- `bg-surface` / `bg-surface-high` / `bg-surface-highest` — neutral surface ladder
- `bg-coral` — primary action
- `text-ink` — primary text
- `bg-success-container`, `bg-warning-container`, `bg-info-container`, `bg-error-container` — feedback surfaces

### Stroke / border colors

- `--color-stroke-soft` — default cards, inputs, shell boundaries
- `--color-stroke-strong` — emphasized separators and focused data surfaces
- `--color-stroke-inverse` — light ring on raised cards

### Radius

- `rounded-card` — default card shell
- `rounded-panel` — raised panels and trays
- `rounded-hero` — hero blocks and marketing-sized surfaces
- `rounded-well` — compact feedback and toggle wells

### Spacing

- `--spacing-gutter` — page gutters
- `--spacing-section` — vertical spacing between major sections
- `--spacing-card` — internal card padding rhythm
- `--spacing-stack` / `--spacing-cluster` — content stacks and chip rows

### Shadows

- `shadow-overshoot` — cards, trays, icon wells
- `shadow-sticker` — coral primary actions and emphasized promotional surfaces

### Typography scale

- `--text-hero` — hero statements
- `--text-display` — large section-level numbers and displays
- `--text-title` — component titles
- `--text-body` — default body copy
- `--text-caption` — labels and metadata

### Container widths

- `--container-reading` — long-form readable content
- `--container-content` — main application canvas
- `--container-shell` — full authenticated layout shell

## 2. Responsive breakpoint usage

Breakpoints are documented as:

- `--breakpoint-tablet` (`48rem`) — switch from stacked cards to simple two-column layouts
- `--breakpoint-desktop` (`64rem`) — enable authenticated shell structure: top bar + fixed sidebar + content canvas
- `--breakpoint-wide` (`80rem`) — allow wider comparison layouts and denser data surfaces

### Usage rules

1. Start mobile-first.
2. Introduce the desktop sidebar only at `lg`/desktop and above.
3. Keep hero and foundation sections stacked until tablet unless side-by-side comparison materially improves comprehension.
4. Data tables may densify at desktop, but empty, error, and action states must remain legible on mobile.
5. Do NOT invent screen-specific breakpoint values without first promoting them into tokens or documenting the exception.

## 3. Foundation rules

### Page backgrounds

- Use `bg-page` for the root canvas.
- Build upward with `surface -> high -> highest` rather than arbitrary grays.

### Cards and surfaces

- Default to `Card` for shared surface behavior.
- Prefer large radii (`rounded-card`, `rounded-panel`) and shared shadow tokens.

### Borders and outlines

- Default to soft outline.
- Use strong outline only for emphasized states, selected shells, or stronger data grouping.

### Icon style

- Icons live in rounded white wells with the shared overshoot shadow.
- Avoid raw icons floating directly on busy surfaces when used as status markers.

### Section headers

- Uppercase micro-label with tracking for metadata.
- Heavy title.
- Softer body copy below the title.

### Button hierarchy

- Primary = coral, deepest emphasis.
- Secondary = neutral raised surface.
- Ghost = quiet action for low-emphasis affordances.

### Inputs and textarea styling

- Rounded 2xl shape, soft outline, white surface.
- Coral focus treatment.
- Error and success messaging stay inside the same shared field primitive.

### Badges / tags

- Use semantic tones (`default`, `success`, `info`, `warning`, `error`, `muted`).
- Keep pills compact and uppercase.

### Toggles / sliders

- Toggles use clear switch semantics with green active track.
- Sliders may use small CSS enhancements only for thumb styling.

### Tables / data list surfaces

- Header row sits on a neutral surface.
- Rows stay on white with soft separators.
- Use `AssetRow` for richer media/data-list rows where needed.

### Empty states

- Use one icon well, one clear title, body copy, and one next action.
- Avoid dead-end messaging.

### Feedback states

- Success = green container
- Info = blue container
- Warning = warm neutral/orange container
- Error = red container

### Toast / alert styling

- `Toast` = non-blocking status update (`role=status`)
- `Alert` = corrective/blocking state (`role=alert`)
- Both share the same container language and action slot

## 4. Shared base components

Implemented under `src/shared/components/ui`:

- `button.tsx`
- `field.tsx`
- `card.tsx`
- `badge.tsx`
- `toggle.tsx`
- `toast.tsx`
- `topbar.tsx`
- `sidebar-nav.tsx`
- `data-table.tsx`
- `metric-progress-card.tsx`
- `empty-state.tsx`

## 5. Composition rule

New screens should compose these primitives FIRST. If a requirement cannot be expressed with the current layer, extend the shared primitive or token set before adding feature-local styling.
