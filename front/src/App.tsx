import {
  CircleCheck,
  Clapperboard,
  CreditCard,
  Dot,
  LayoutDashboard,
  Library,
  PencilLine,
  Plus,
  Sparkles,
} from 'lucide-react'
import { ActionChip } from './shared/components/ui/action-chip'
import { Badge } from './shared/components/ui/badge'
import { Button } from './shared/components/ui/button'
import { Card, SectionShell } from './shared/components/ui/card'
import { AssetRow, DataTable } from './shared/components/ui/data-table'
import { EmptyState } from './shared/components/ui/empty-state'
import { Input, PromptField } from './shared/components/ui/field'
import { MetricProgressCard } from './shared/components/ui/metric-progress-card'
import { SidebarNav } from './shared/components/ui/sidebar-nav'
import { Slider } from './shared/components/ui/slider'
import { Alert, Toast } from './shared/components/ui/toast'
import { Toggle } from './shared/components/ui/toggle'
import { Topbar } from './shared/components/ui/topbar'

const NAV_ITEMS = [
  {
    href: '/',
    label: 'Dashboard',
    shortLabel: 'Home',
    active: true,
    icon: LayoutDashboard,
  },
  { href: '/', label: 'Studio', shortLabel: 'Studio', icon: Clapperboard },
  { href: '/', label: 'Library', shortLabel: 'Library', icon: Library },
  { href: '/', label: 'Billing', shortLabel: 'Billing', icon: CreditCard },
]

const ACTION_CHIPS = ['Prompt ready', 'Auto save', 'Export clean']
const META_TAGS = ['Animation', 'Storyboard']

const TABLE_COLUMNS = [
  { key: 'asset', header: 'Asset' },
  { key: 'owner', header: 'Owner' },
  { key: 'status', header: 'Status' },
  { key: 'updated', header: 'Updated', align: 'right' },
] as const

const TABLE_ROWS = [
  {
    id: 'scene-01',
    asset: (
      <AssetRow
        label="Launch storyboard"
        meta="16 scenes • 4K thumbnail pack"
        status="Ready"
      />
    ),
    owner: 'Creative ops',
    status: <Badge tone="success">Ready</Badge>,
    updated: '2m ago',
  },
  {
    id: 'scene-02',
    asset: (
      <AssetRow
        label="Voice timing pass"
        meta="Narration alignment • 8 checkpoints"
        status="Processing"
      />
    ),
    owner: 'Narrative',
    status: <Badge tone="info">Processing</Badge>,
    updated: '14m ago',
  },
  {
    id: 'scene-03',
    asset: (
      <AssetRow
        label="Thumbnail contrast review"
        meta="Poster frame QA • Brand safety"
        status="Needs review"
      />
    ),
    owner: 'Brand',
    status: <Badge tone="warning">Needs review</Badge>,
    updated: '1h ago',
  },
]

const ICON_CATALOG = [
  { label: 'Sparkles', Icon: Sparkles },
  { label: 'Pencil line', Icon: PencilLine },
  { label: 'Circle check', Icon: CircleCheck },
  { label: 'Dashboard', Icon: LayoutDashboard },
  { label: 'Clapperboard', Icon: Clapperboard },
  { label: 'Library', Icon: Library },
  { label: 'Credit card', Icon: CreditCard },
  { label: 'Plus', Icon: Plus },
]

export function App() {
  return (
    <div className="foundation-page">
      <div className="mx-auto flex min-h-screen max-w-[var(--container-shell)]">
        <SidebarNav
          className="fixed inset-y-0 left-0"
          footer={
            <div className="space-y-3">
              <p className="text-[10px] font-extrabold uppercase tracking-section text-white/75">
                Palette seed
              </p>
              <p className="text-xl font-black tracking-tight">
                Strict tokens, calm neutrals, and one reusable system.
              </p>
            </div>
          }
          items={NAV_ITEMS}
          title="Foundations"
        />

        <div className="w-full lg:ml-72">
          <Topbar
            actions={
              <>
                <ActionChip
                  icon={
                    <Dot
                      aria-hidden="true"
                      className="size-4"
                      strokeWidth={2.6}
                    />
                  }
                >
                  Live preview
                </ActionChip>
                <Button>Publish kit</Button>
              </>
            }
            description="Responsive shells, tokens, and primitives now share one source of truth."
            eyebrow="Top bar"
            title="Design showcase canvas"
          />

          <main className="mx-auto flex max-w-[var(--container-content)] flex-col gap-8 px-4 pb-28 pt-6 sm:px-6 lg:px-10 lg:pb-12 lg:pt-10">
            <SectionShell
              description="A token-first SaaS foundation that mirrors the requested reference style: bold headline weight, coral actions, pill controls, soft container layering, and reusable responsive rules."
              eyebrow="Hero title"
              title="Foundational components and palette for a calm but high-contrast creative SaaS."
            >
              <Card className="grid gap-6 rounded-hero bg-white p-6 sm:p-8 lg:grid-cols-[1.1fr_0.9fr] lg:p-10">
                <div className="space-y-6">
                  <div className="flex flex-wrap gap-2">
                    <Badge>Foundation</Badge>
                    <Badge tone="success">Visual contract locked</Badge>
                  </div>
                  <div className="space-y-4">
                    <h2 className="max-w-3xl text-4xl font-black leading-[0.95] tracking-[-0.06em] text-ink sm:text-5xl lg:text-[var(--text-hero)]">
                      Build once. Keep every future screen visually consistent.
                    </h2>
                    <p className="max-w-2xl text-base font-medium leading-7 text-ink/62">
                      Tokens now cover color, spacing, type scale, borders,
                      shadows, containers, and documented breakpoint usage so
                      later feature work does NOT drift.
                    </p>
                  </div>
                  <div className="flex flex-wrap gap-3">
                    <Button>Generate</Button>
                    <Button variant="secondary">Save draft</Button>
                    <Button variant="ghost">Preview only</Button>
                    <Button isLoading>Generating</Button>
                  </div>
                  <div className="flex flex-wrap gap-3">
                    {ACTION_CHIPS.map((chip) => (
                      <ActionChip
                        icon={
                          <Dot
                            aria-hidden="true"
                            className="size-4"
                            strokeWidth={2.6}
                          />
                        }
                        key={chip}
                      >
                        {chip}
                      </ActionChip>
                    ))}
                  </div>
                </div>

                <div className="grid gap-4">
                  <MetricProgressCard
                    helper="Desktop shells stay anchored while mobile collapses to top bar, bottom nav, and a single floating action."
                    label="Responsive readiness"
                    progress={84}
                    trend="+18%"
                    value="84%"
                  />
                  <Toast
                    action={
                      <ActionChip
                        icon={
                          <Dot
                            aria-hidden="true"
                            className="size-4"
                            strokeWidth={2.6}
                          />
                        }
                      >
                        Open docs
                      </ActionChip>
                    }
                    icon={
                      <Sparkles
                        aria-hidden="true"
                        className="size-4"
                        strokeWidth={2}
                      />
                    }
                    title="System status synced"
                    tone="success"
                  >
                    Guidance now documents breakpoints, surfaces, and component
                    hierarchy in reusable form.
                  </Toast>
                </div>
              </Card>
            </SectionShell>

            <SectionShell
              description="Foundation rules are explicit instead of implied: backgrounds, surfaces, borders, icons, section headers, controls, feedback, data surfaces, and empty states all point back to the same tokens."
              eyebrow="Foundation Rules"
              title="Surface, border, and feedback rules are part of the system now."
            >
              <div className="grid gap-5 lg:grid-cols-3">
                <Card className="space-y-4 bg-surface">
                  <p className="foundation-section-eyebrow">Page backgrounds</p>
                  <p className="text-lg font-black tracking-tight">
                    Page → surface → high → highest
                  </p>
                  <p className="foundation-body">
                    Use `bg-page` for the canvas, then climb the surface ladder
                    for cards and focused trays.
                  </p>
                  <div className="grid gap-3 sm:grid-cols-3">
                    {[
                      ['Page', 'bg-page'],
                      ['Surface', 'bg-surface'],
                      ['High', 'bg-surface-high'],
                    ].map(([label, className]) => (
                      <div className="space-y-2" key={label}>
                        <div
                          className={`h-16 rounded-swatch ${className} shadow-overshoot`}
                        />
                        <p className="text-xs font-bold text-ink/58">{label}</p>
                      </div>
                    ))}
                  </div>
                </Card>

                <Card className="space-y-4 bg-white">
                  <p className="foundation-section-eyebrow">
                    Borders and icon style
                  </p>
                  <div className="flex items-center gap-3">
                    <span className="foundation-icon">
                      <Sparkles
                        aria-hidden="true"
                        className="size-4"
                        strokeWidth={2}
                      />
                    </span>
                    <div>
                      <p className="font-bold text-ink">Rounded icon wells</p>
                      <p className="text-sm font-medium text-ink/52">
                        Icons live inside soft white wells with one outline and
                        shared depth.
                      </p>
                    </div>
                  </div>
                  <div className="rounded-panel bg-surface-high p-4 foundation-outline-soft">
                    <p className="text-sm font-bold text-ink">Soft outline</p>
                    <p className="text-sm font-medium text-ink/52">
                      Default card and control boundary.
                    </p>
                  </div>
                  <div className="rounded-panel bg-white p-4 foundation-outline-strong">
                    <p className="text-sm font-bold text-ink">Strong outline</p>
                    <p className="text-sm font-medium text-ink/52">
                      Reserved for emphasized states and table headers.
                    </p>
                  </div>
                </Card>

                <Card className="space-y-4 bg-surface-highest">
                  <p className="foundation-section-eyebrow">
                    Section headers and hierarchy
                  </p>
                  <p className="foundation-section-title">
                    Uppercase metadata, heavy titles, softer body copy.
                  </p>
                  <p className="foundation-body">
                    Hero copy stays dominant, section titles remain compact, and
                    support text never competes with the action hierarchy.
                  </p>
                  <div className="flex flex-wrap gap-3">
                    <Button size="lg">Primary action</Button>
                    <Button variant="secondary">Secondary button</Button>
                    <Button variant="ghost">Ghost button</Button>
                  </div>
                </Card>
              </div>
            </SectionShell>

            <SectionShell
              description="Buttons, fields, tags, sliders, toggles, and feedback callouts now behave like a coherent kit instead of isolated examples."
              eyebrow="Base Components"
              title="The first reusable components are implemented and ready to compose."
            >
              <div className="grid gap-5 lg:grid-cols-[1.05fr_0.95fr]">
                <Card className="space-y-5 bg-surface-high">
                  <div className="flex flex-wrap gap-3">
                    <Button>Generate</Button>
                    <Button variant="secondary">Save draft</Button>
                    <Button variant="ghost">Preview only</Button>
                    <Button isLoading>Generating</Button>
                  </div>
                  <div className="flex flex-wrap gap-2.5">
                    <Badge variant="dark">#vector</Badge>
                    {META_TAGS.map((tag) => (
                      <Badge key={tag} variant="outline">
                        {tag}
                      </Badge>
                    ))}
                    <Badge variant="removable">Coral tag</Badge>
                  </div>
                  <div className="grid gap-4 lg:grid-cols-2">
                    <Input
                      defaultValue="Brandtoon launch kit"
                      hint="Use sentence case and concise naming."
                      label="Workspace name"
                      message="Required field"
                      state="error"
                    />
                    <Input
                      defaultValue="UI system, tokens, hero polish"
                      hint="Comma-separated tags help surface related flows later."
                      label="Tags"
                      message="Looks good"
                      state="success"
                    />
                  </div>
                  <PromptField
                    icon={
                      <PencilLine
                        aria-hidden="true"
                        className="size-4"
                        strokeWidth={2}
                      />
                    }
                    placeholder="Describe the visual goal"
                    title="Prompt"
                  />
                </Card>

                <div className="grid gap-5">
                  <Card className="space-y-4 bg-white">
                    <Slider
                      defaultValue={62}
                      label="Intensity"
                      markers={[
                        { label: 'minimal', value: 'subtle' },
                        { label: 'playful', value: 'balanced' },
                        { label: 'chaotic', value: 'maxed' },
                      ]}
                      max={100}
                      min={0}
                      valueLabel="VIBRANT"
                    />
                    <Toggle
                      defaultChecked
                      description="Boost the palette and illustration energy automatically."
                      label="Vibrant rendering"
                    />
                    <Toggle
                      description="Keep a quieter treatment for supporting accents."
                      label="Subtle supporting details"
                    />
                  </Card>

                  <Toast
                    dismissLabel="Close toast"
                    icon={
                      <CircleCheck
                        aria-hidden="true"
                        className="size-4"
                        strokeWidth={2}
                      />
                    }
                    title="Project saved"
                    tone="success"
                  >
                    Your latest prompt and settings were synced successfully.
                  </Toast>

                  <Alert
                    action={
                      <ActionChip
                        icon={
                          <Dot
                            aria-hidden="true"
                            className="size-4"
                            strokeWidth={2.6}
                          />
                        }
                      >
                        Retry
                      </ActionChip>
                    }
                    icon={
                      <Sparkles
                        aria-hidden="true"
                        className="size-4"
                        strokeWidth={2}
                      />
                    }
                    title="Contrast exception"
                    tone="error"
                  >
                    Use alerts for blocking or corrective states, not decorative
                    accents.
                  </Alert>

                  <Card className="space-y-4 bg-surface-high">
                    <div className="space-y-1">
                      <p className="foundation-section-eyebrow">Icon catalog</p>
                      <p className="text-sm font-bold text-ink">
                        Lucide is now wired as the third-party icon library.
                      </p>
                    </div>
                    <div className="grid grid-cols-2 gap-3 sm:grid-cols-3">
                      {ICON_CATALOG.map(({ label, Icon }) => (
                        <div
                          className="rounded-panel bg-white p-3 foundation-outline-soft"
                          key={label}
                        >
                          <div className="mb-2 inline-flex rounded-xl bg-surface p-2 text-ink">
                            <Icon aria-hidden="true" className="size-4" />
                          </div>
                          <p className="text-xs font-bold text-ink/62">
                            {label}
                          </p>
                        </div>
                      ))}
                    </div>
                  </Card>
                </div>
              </div>
            </SectionShell>

            <SectionShell
              description="Navigation, table surfaces, empty states, and progress cards are now implemented as shared building blocks for future product screens."
              eyebrow="Data & Status Foundations"
              title="App shell primitives and data surfaces are ready for feature composition."
            >
              <div className="grid gap-5 xl:grid-cols-[1.15fr_0.85fr]">
                <div className="space-y-5">
                  <DataTable columns={TABLE_COLUMNS} rows={TABLE_ROWS} />
                  <DataTable
                    columns={TABLE_COLUMNS}
                    emptyDescription="Create your first asset row to start using the shared data surface."
                    emptyTitle="No assets yet"
                    rows={[]}
                  />
                </div>

                <div className="grid gap-5">
                  <MetricProgressCard
                    helper="Compared with the previous export cycle."
                    label="Completion"
                    progress={72}
                    trend="+12%"
                    value="72%"
                  />
                  <EmptyState
                    action={
                      <Button variant="secondary">Create first asset</Button>
                    }
                    icon={
                      <Sparkles
                        aria-hidden="true"
                        className="size-4"
                        strokeWidth={2}
                      />
                    }
                    title="Empty states stay optimistic, not dead"
                  >
                    Use one icon well, a clear next action, and body copy that
                    explains what becomes available after the first item exists.
                  </EmptyState>
                </div>
              </div>
            </SectionShell>
          </main>

          <button
            className="fixed bottom-24 right-5 z-20 inline-flex size-14 items-center justify-center rounded-full bg-coral text-2xl font-black text-white shadow-sticker lg:hidden"
            type="button"
          >
            <Plus aria-hidden="true" className="size-6" strokeWidth={2.6} />
          </button>

          <nav
            aria-label="Mobile bottom navigation"
            className="fixed inset-x-0 bottom-0 z-20 flex items-center justify-around border-t border-[color:var(--color-stroke-soft)] bg-white/92 px-3 py-3 backdrop-blur-xl lg:hidden"
          >
            {NAV_ITEMS.map((item) => {
              const Icon = item.icon

              return (
                <a
                  className="flex min-w-16 flex-col items-center gap-1 text-[11px] font-bold text-ink/62"
                  href={item.href}
                  key={item.shortLabel}
                >
                  <span className={item.active ? 'text-coral' : 'text-current'}>
                    <Icon
                      aria-hidden="true"
                      className="size-4"
                      strokeWidth={2.4}
                    />
                  </span>
                  <span>{item.shortLabel}</span>
                </a>
              )
            })}
          </nav>
        </div>
      </div>
    </div>
  )
}
