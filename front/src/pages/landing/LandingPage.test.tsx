import { fireEvent, render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { LandingPage } from './LandingPage'

const navigateMock = vi.fn()

vi.mock('lucide-react', () => ({
  ArrowRight: () => null,
  CheckCircle2: () => null,
  Film: () => null,
  Play: () => null,
  Sparkles: () => null,
}))

vi.mock('../../shared/components/ui/action-chip', () => ({
  ActionChip: ({ children }: { children: unknown }) => <div>{children}</div>,
}))

vi.mock('../../shared/components/ui/badge', () => ({
  Badge: ({ children }: { children: unknown }) => <div>{children}</div>,
}))

vi.mock('../../shared/components/ui/button', () => ({
  Button: ({
    children,
    onClick,
  }: { children: unknown; onClick?: () => void }) => (
    <button onClick={onClick} type="button">
      {children}
    </button>
  ),
}))

vi.mock('../../shared/components/ui/card', () => ({
  Card: ({ children }: { children: unknown }) => <div>{children}</div>,
  SectionShell: ({ children }: { children: unknown }) => (
    <section>{children}</section>
  ),
}))

vi.mock('../../shared/components/ui/field', () => ({
  Input: () => <input aria-label="Mocked landing input" />,
}))

vi.mock('../../shared/components/ui/slider', () => ({
  Slider: () => <div>Mocked slider</div>,
}))

vi.mock('../../shared/components/ui/topbar', () => ({
  Topbar: ({ actions, title }: { actions?: unknown; title?: string }) => (
    <header>
      <span>{title}</span>
      {actions}
    </header>
  ),
}))

vi.mock('react-router-dom', () => ({
  useNavigate: () => navigateMock,
}))

describe('LandingPage', () => {
  beforeEach(() => {
    navigateMock.mockReset()
  })

  it('navigates to login with the creative studio redirect when the hero CTA is clicked', () => {
    render(<LandingPage />)

    fireEvent.click(screen.getByRole('button', { name: /create your avatar/i }))

    expect(navigateMock).toHaveBeenCalledWith('/login?next=%2Fcreative-studio')
  })
})
