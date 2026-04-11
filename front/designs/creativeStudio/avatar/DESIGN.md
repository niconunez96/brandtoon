# Kinetic Ink

### 1. Overview & Creative North Star
**Creative North Star: The Energetic Editorial**
Kinetic Ink is a design system that marries high-authority typography with the playful, tactile energy of physical sticker art. It moves away from the "clinical SaaS" look by embracing heavy ink-weights, high-contrast color pops, and a "thick" architectural feel. The system prioritizes bold intent over subtle utility, using intentional asymmetry and oversized elements to create a rhythmic, curated experience.

### 2. Colors
The palette is built on a foundation of "Dark Neutral" ink (#2D3436) and "Energetic Coral" (#FF6B6B). 

*   **The "No-Line" Rule:** Sectioning is achieved through shifts in surface containers (e.g., transitioning from `surface` #fbfcfc to `surface_container` #f1f2f6). 1px solid borders are strictly prohibited for structural layout; use them only as "ink-lines" for specific component definitions at 10% opacity.
*   **Surface Hierarchy:** 
    *   `surface_container_lowest`: Pure white cards (#ffffff).
    *   `surface_container`: Page foundations (#f1f2f6).
    *   `surface_container_highest`: Active/Hover states or sidebar depths (#dfe6e9).
*   **Glass & Gradient:** Floating navigation bars must utilize 80% opacity with 24px backdrop blurs. 
*   **Signature Textures:** Primary CTAs use a "thick ink" approach, often paired with a physical-depth shadow.

### 3. Typography
Kinetic Ink uses **Plus Jakarta Sans** across all levels to maintain a modern, geometric clarity. The scale is intentionally dramatic.

*   **Display Large (3.75rem / 60px):** Black weight (900). Use for core value propositions.
*   **Headline Medium (1.875rem / 30px):** ExtraBold. For section starts.
*   **Body Large (1.125rem / 18px):** Medium weight. Optimized for editorial readability with a 1.6x line height.
*   **Labels (10px - 11px):** Black weight (900) with 20% letter-spacing (0.2em). Always uppercase. This "micro-labeling" provides the architectural/technical contrast to the playful headlines.

### 4. Elevation & Depth
Depth in Kinetic Ink is tactile rather than atmospheric.

*   **The Sticker Effect:** Cards utilize a dual shadow:
    1.  An "Overshoot" shadow: `0 8px 24px -4px rgba(45, 52, 54, 0.12)`
    2.  A "Hard Base": `0 4px 0 0 rgba(45, 52, 54, 0.05)`
*   **The Layering Principle:** Stacking follows a logic of `surface` -> `surface_container` -> `surface_container_lowest` (the "Raised Card").
*   **Primary Action Shadow:** Specific to buttons, use a heavy-tinted offset: `0 8px 0 -2px rgba(170, 44, 50, 0.2), 0 12px 24px -4px rgba(170, 44, 50, 0.3)`.

### 5. Components
*   **Buttons:** Must be `rounded-2xl` (1rem). Primary buttons utilize the "Energetic Coral" container with the Primary Action Shadow.
*   **Sticker Cards:** `rounded-[2.5rem]` (40px) containers. These represent the highest level of brand expression, often featuring 2px borders at 10% opacity.
*   **Micro-Tags:** Pill-shaped (`rounded-full`) with uppercase 10px Black weight text. High contrast backgrounds only.
*   **Navigation:** Sidebars should use `surface_bright` with `surface_container` active states, emphasizing horizontal "slide-in" hover effects.

### 6. Do's and Don'ts
*   **Do:** Use extreme font weight contrasts (Black 900 vs Medium 500).
*   **Do:** Embrace large corner radii (up to 40px) for container-level elements.
*   **Don't:** Use generic 1px #E2E8F0 borders. Use color shifts or 10% opacity ink lines.
*   **Don't:** Use low-contrast text. Ensure `on_surface_variant` remains highly legible (#2D3436).
*   **Do:** Use "Sticker Shadows" to make important UI elements feel like they can be peeled off the screen.