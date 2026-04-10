import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes, useLocation } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '../../services/auth.api'
import { reloadBrowserWindow } from '../../shared/lib/browser'
import { ProtectedRoute } from './ProtectedRoute'

const useCurrentUserQueryMock = vi.fn()

vi.mock('../../queries/useCurrentUserQuery', () => ({
  useCurrentUserQuery: () => useCurrentUserQueryMock(),
}))

vi.mock('../../shared/lib/browser', async () => {
  const actual = await vi.importActual<
    typeof import('../../shared/lib/browser')
  >('../../shared/lib/browser')

  return {
    ...actual,
    reloadBrowserWindow: vi.fn(),
  }
})

function LocationProbe() {
  const location = useLocation()

  return (
    <div data-testid="location-display">{`${location.pathname}${location.search}`}</div>
  )
}

function renderProtectedRoute(initialEntry = '/creative-studio?tab=assets') {
  return render(
    <MemoryRouter initialEntries={[initialEntry]}>
      <Routes>
        <Route path="/login" element={<div>Login page</div>} />
        <Route element={<ProtectedRoute />}>
          <Route path="/creative-studio" element={<div>Creative Studio</div>} />
        </Route>
      </Routes>
      <LocationProbe />
    </MemoryRouter>,
  )
}

describe('ProtectedRoute', () => {
  beforeEach(() => {
    useCurrentUserQueryMock.mockReset()
    vi.mocked(reloadBrowserWindow).mockReset()
  })

  it('renders the loading state while the session query is pending', () => {
    useCurrentUserQueryMock.mockReturnValue({
      error: null,
      isError: false,
      isLoading: true,
    })

    renderProtectedRoute()

    expect(screen.getByText(/loading session/i)).toBeInTheDocument()
    expect(screen.getByText(/preparing your studio/i)).toBeInTheDocument()
  })

  it('redirects unauthorized visitors to login and preserves the next destination', async () => {
    useCurrentUserQueryMock.mockReturnValue({
      error: new ApiError('Unauthorized', 401),
      isError: true,
      isLoading: false,
    })

    renderProtectedRoute('/creative-studio?tab=assets&filter=all')

    expect(await screen.findByText(/login page/i)).toBeInTheDocument()
    expect(screen.getByTestId('location-display')).toHaveTextContent(
      '/login?next=%2Fcreative-studio%3Ftab%3Dassets%26filter%3Dall',
    )
  })

  it('renders a retry state for non-unauthorized session failures', async () => {
    const user = userEvent.setup()
    const reloadBrowserWindowMock = vi.mocked(reloadBrowserWindow)

    useCurrentUserQueryMock.mockReturnValue({
      error: new ApiError('Server error', 500),
      isError: true,
      isLoading: false,
    })

    renderProtectedRoute()

    expect(
      screen.getByText(/we could not verify your session/i),
    ).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: /try again/i }))

    expect(reloadBrowserWindowMock).toHaveBeenCalledTimes(1)
  })

  it('renders the protected outlet when the session is available', async () => {
    useCurrentUserQueryMock.mockReturnValue({
      data: {
        user: {
          avatarUrl: 'https://avatar.example.com/nico.png',
          email: 'nico@example.com',
          id: 'user-1',
          name: 'Nico',
        },
      },
      error: null,
      isError: false,
      isLoading: false,
    })

    renderProtectedRoute()

    expect(await screen.findByText(/creative studio/i)).toBeInTheDocument()
    expect(screen.getByTestId('location-display')).toHaveTextContent(
      '/creative-studio?tab=assets',
    )
  })
})
