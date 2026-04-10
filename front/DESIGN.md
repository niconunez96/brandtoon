# Design System Strategy: The Tactile Pop Editorial

## 1. Overview & Creative North Star
This design system is defined by a "Creative North Star" we call **The Tactile Pop Editorial**.

The objective is to dismantle the rigidity of standard enterprise SaaS and replace it with a layout that feels curated, high-end, and undeniably energetic. We are blending the sophisticated structure of a premium fashion editorial with the approachable, "sticker-like" physics of modern comic aesthetics. By utilizing intentional asymmetry, overlapping containers, and a "bouncy" interactive language, we move away from "templates" and toward a signature digital experience.

## 2. Colors & Surface Logic
The palette centers around a high-vibrancy primary coral, supported by a sophisticated range of tinted neutrals and pastel containers.

### The "No-Line" Rule
Standard 1px grey borders are strictly prohibited for sectioning. Layout boundaries must be defined through:

- **Tonal Shifts:** Moving from `surface` to `surface-container-low`.
- **Vibrant Pastels:** Utilizing `secondary-container` (`#ffc2cc`) or `tertiary-container` (`#ca9cff`) for card backgrounds to create a "sticker" contrast against the `background` (`#fff4f4`).

### Surface Hierarchy & Nesting
Treat the UI as a physical desk of layered paper. Use the surface tiers to establish importance:

1. **Base Layer:** `surface` (`#fff4f4`)
2. **Sectioning:** `surface-container-low` (`#ffecee`)
3. **Floating Elements (Cards):** `surface-container-lowest` (`#ffffff`) or vibrant pastel containers.
4. **Interactive Overlays:** Use `surface-bright` with a **Glassmorphism effect** (`backdrop-blur: 12px; opacity: 80%`) to allow background colors to bleed through, softening the "sticker" edges.

### Signature Textures
Main CTAs and Hero sections should avoid flat fills. Instead, apply a **Linear Gradient** (45deg) transitioning from `primary` (`#aa2c32`) to `primary-container` (`#ff7574`). This adds "soul" and depth that prevents the "Playful Professional" style from looking juvenile.

## 3. Typography: The Character Scale
We use **Plus Jakarta Sans** as our sole typeface, but we leverage weight to inject personality.

- **Display & Headlines:** Must use **ExtraBold (800)**. This weight provides the "comic" character needed to balance the soft coral tones. Use tighter letter-spacing (`-0.02em`) for a high-end editorial feel.
- **Titles:** Use **Bold (700)** to maintain authority without over-powering the body content.
- **Body:** Keep to **Medium (500)** or **Regular (400)**. Our `body-lg` (`1rem`) is the workhorse.
- **Labels:** Always **SemiBold (600)** and uppercase when used in buttons or navigation to ensure they feel like "stamps."

## 4. Elevation & Depth: The Sticker Effect
Unlike traditional Material Design which uses high-diffusion ambient shadows, this system uses a "Pop" shadow logic.

### The Stacking Principle

- **Shadow Construction:** Use a low-diffusion (Blur: `4px`–`8px`), slightly offset (`Y: 4px`) shadow. The shadow color must never be grey; use a 15% opacity version of `on-surface` (`#4c212b`) to create a warm, natural "lift."
- **The Comic Outline:** To achieve the "Playful Professional" look, every card and primary button must have a **1.5px border**. The border color should be a darker tonal variant of the background it sits on (e.g., if a card is `secondary-container`, the border is `on-secondary-container` at 20% opacity).
- **Corner Radii:** Strictly adhere to the **20px-24px (md/lg)** range for all primary containers. This "super-ellipse" feel makes the UI feel soft and touchable.

## 5. Components

### Buttons

- **Primary:** Gradient fill (`primary` to `primary-container`), 24px radius, with a 2px "darker coral" border.
- **Micro-interaction:** On hover, the button should scale (`1.05`) and the shadow should increase in offset. Use a `cubic-bezier(0.34, 1.56, 0.64, 1)` transition for a "bounce" feel.

### Cards & Containers

- **Rule:** Forbid the use of divider lines.
- **Separation:** Content within cards must be separated by vertical whitespace (utilizing the `1.5rem` spacing token) or by nesting a `surface-container-high` inner box inside a `surface-container-low` outer box.

### Input Fields

- **Style:** Use `surface-container-lowest` as the fill.
- **Borders:** A 2px border using `outline-variant` (`#dc9ca8`). On focus, the border shifts to `primary` (`#aa2c32`) and the entire field "pops" (slight scale increase).

### Chips & Badges

- **Style:** Use `tertiary-container` for selection. Use `full` (`9999px`) roundness to contrast against the 24px cards.

## 6. Do's and Don'ts

### Do

- **DO** use intentional asymmetry. Overlap a card slightly over a section header to break the grid.
- **DO** use "Bounce" easing for all transitions. A "Playful Professional" never slides; they pop.
- **DO** use the `primary-fixed` and `secondary-fixed` tokens for elements that need to remain vibrant regardless of dark/light mode shifts.

### Don't

- **DON'T** use pure black (`#000000`) for text or shadows. It kills the "Energetic Coral" warmth. Use `on-surface` (`#4c212b`).
- **DON'T** use 1px thin lines to separate list items. Use `8px` or `12px` of vertical gap.
- **DON'T** use sharp corners. Anything less than 12px radius is a violation of the system's silhouette.
