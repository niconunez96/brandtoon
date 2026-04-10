import { render, screen } from '@testing-library/react'
import { DataTable } from './data-table'
import { Input, PromptField, Textarea } from './field'
import { MetricProgressCard } from './metric-progress-card'
import { SidebarNav } from './sidebar-nav'
import { Alert, Toast } from './toast'
import { Toggle } from './toggle'
import { Topbar } from './topbar'

describe('foundation primitives', () => {
  it('renders navigation shell primitives with accessible landmarks', () => {
    render(
      <>
        <Topbar
          actions={<button type="button">Publish</button>}
          eyebrow="Workspace"
          title="Design foundations"
        />
        <SidebarNav
          items={[
            { href: '/', label: 'Dashboard', shortLabel: 'Home' },
            {
              href: '/studio',
              label: 'Studio',
              shortLabel: 'Studio',
              active: true,
            },
          ]}
          title="Brandtoon UI"
        />
      </>,
    )

    expect(screen.getByRole('banner')).toBeInTheDocument()
    expect(screen.getByLabelText('Primary sidebar')).toBeInTheDocument()
    expect(screen.getByText('Design foundations')).toBeInTheDocument()
    expect(screen.getByText('Studio')).toBeInTheDocument()
  })

  it('renders data surfaces, empty states, and progress metrics', () => {
    const columns = [
      { key: 'asset', header: 'Asset' },
      { key: 'status', header: 'Status' },
    ] as const

    render(
      <>
        <DataTable
          columns={columns}
          emptyDescription="Create your first asset to populate the queue."
          emptyTitle="No assets yet"
          rows={[]}
        />
        <DataTable
          columns={columns}
          rows={[
            { id: 'scene-01', asset: 'Launch storyboard', status: 'Ready' },
          ]}
        />
        <MetricProgressCard
          helper="Compared with the previous export cycle."
          label="Completion"
          progress={72}
          trend="+12%"
          value="72%"
        />
      </>,
    )

    expect(screen.getByText('No assets yet')).toBeInTheDocument()
    expect(screen.getByText('Launch storyboard')).toBeInTheDocument()
    expect(
      screen.getByRole('progressbar', { name: 'Completion' }),
    ).toHaveAttribute('aria-valuenow', '72')
  })

  it('supports alert states, field validation, and interactive toggles', () => {
    render(
      <>
        <Alert title="Render failed" tone="error">
          Retry with a shorter prompt.
        </Alert>
        <Input
          hint="Name should stay under 40 characters."
          label="Workspace name"
          message="Required field"
          state="error"
        />
        <Textarea
          hint="Keep the message concise."
          label="Creative brief"
          message="Looks good"
          state="success"
        />
        <PromptField placeholder="Describe the visual goal" title="Prompt" />
        <Toggle
          defaultChecked
          description="Use the stronger coral hierarchy."
          label="Design mode"
        />
        <Toast title="Project saved">Your latest prompt was synced.</Toast>
      </>,
    )

    expect(screen.getByRole('alert')).toBeInTheDocument()
    expect(screen.getByRole('status')).toBeInTheDocument()
    expect(screen.getByText('Required field')).toBeInTheDocument()
    expect(screen.getByText('Looks good')).toBeInTheDocument()
    expect(screen.getByText('0 / 500')).toBeInTheDocument()
    expect(screen.getByRole('switch', { name: 'Design mode' })).toHaveAttribute(
      'aria-checked',
      'true',
    )
  })
})
