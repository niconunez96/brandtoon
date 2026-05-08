import { fireEvent, render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { API_BASE_URL } from '../../shared/config/api'
import { navigateToExternalUrl } from '../../shared/lib/browser'
import { LoginPage } from './LoginPage'

const useCurrentUserQueryMock = vi.fn()
const useSearchParamsMock = vi.fn()

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
}))

vi.mock('../../queries/useCurrentUserQuery', () => ({
  useCurrentUserQuery: () => useCurrentUserQueryMock(),
}))

vi.mock('react-router-dom', () => ({
  Navigate: ({ to }: { to: string }) => <div>Navigate to {to}</div>,
  useSearchParams: () => useSearchParamsMock(),
}))

vi.mock('../../shared/lib/browser', () => ({
  navigateToExternalUrl: vi.fn(),
}))

describe('LoginPage', () => {
  beforeEach(() => {
    useCurrentUserQueryMock.mockReset()
    useSearchParamsMock.mockReset()
    vi.mocked(navigateToExternalUrl).mockReset()
    useCurrentUserQueryMock.mockReturnValue({
      data: undefined,
    })
    useSearchParamsMock.mockReturnValue([
      new URLSearchParams('next=%2Fcreative-studio'),
    ])
  })

  it('sends users to the backend Google auth endpoint with the sanitized redirect', () => {
    render(<LoginPage />)

    fireEvent.click(
      screen.getByRole('button', { name: /continue with google/i }),
    )

    expect(navigateToExternalUrl).toHaveBeenCalledWith(
      `${API_BASE_URL}/auth/google/login?redirectTo=%2Fcreative-studio`,
    )
  })
})
