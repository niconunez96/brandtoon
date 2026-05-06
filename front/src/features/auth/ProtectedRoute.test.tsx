import { fireEvent, render, screen } from '@testing-library/react'
import type { ReactNode } from 'react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '../../services/auth.api'
import { reloadBrowserWindow } from '../../shared/lib/browser'
import { ProtectedRoute } from './ProtectedRoute'

const useCurrentUserQueryMock = vi.fn()
const locationMock = vi.fn()

vi.mock('../../queries/useCurrentUserQuery', () => ({
  useCurrentUserQuery: () => useCurrentUserQueryMock(),
}))

vi.mock('./AvatarGenerationEventsProvider', () => ({
  AvatarGenerationEventsProvider: ({ children }: { children: ReactNode }) => (
    <>{children}</>
  ),
}))

vi.mock('../../shared/lib/browser', async () => {
  const actual = await vi.importActual<typeof import('../../shared/lib/browser')>(
    '../../shared/lib/browser',
  )

  return {
    ...actual,
    reloadBrowserWindow: vi.fn(),
  }
})

vi.mock('react-router-dom', () => ({
  Navigate: ({ to }: { to: string }) => (
    <div data-testid="navigate-target">{to}</div>
  ),
  Outlet: () => <div>Creative Studio</div>,
  useLocation: () => locationMock(),
}))

describe('ProtectedRoute', () => {
  beforeEach(() => {
    useCurrentUserQueryMock.mockReset()
    locationMock.mockReset()
    vi.mocked(reloadBrowserWindow).mockReset()
    locationMock.mockReturnValue({
      pathname: '/creative-studio',
      search: '?tab=assets',
    })
  })

  it('redirects unauthorized visitors to login and preserves the next destination', () => {
    useCurrentUserQueryMock.mockReturnValue({
      error: new ApiError('Unauthorized', 401),
      isError: true,
      isLoading: false,
    })
    locationMock.mockReturnValue({
      pathname: '/creative-studio',
      search: '?tab=assets&filter=all',
    })

    render(<ProtectedRoute />)

    expect(screen.getByTestId('navigate-target')).toHaveTextContent(
      '/login?next=%2Fcreative-studio%3Ftab%3Dassets%26filter%3Dall',
    )
  })

  it('renders a retry state for non-unauthorized session failures', () => {
    useCurrentUserQueryMock.mockReturnValue({
      error: new ApiError('Server error', 500),
      isError: true,
      isLoading: false,
    })

    render(<ProtectedRoute />)

    expect(
      screen.getByText(/we could not verify your session/i),
    ).toBeInTheDocument()

    fireEvent.click(screen.getByRole('button', { name: /try again/i }))

    expect(reloadBrowserWindow).toHaveBeenCalledTimes(1)
  })
})
