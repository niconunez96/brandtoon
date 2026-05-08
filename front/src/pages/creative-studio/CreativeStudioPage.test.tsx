import { fireEvent, render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { CreativeStudioPage } from './CreativeStudioPage'

const navigateMock = vi.fn()
const useAvatarsQueryMock = vi.fn()
const useCreateAvatarMutationMock = vi.fn()
const useCurrentUserQueryMock = vi.fn()

vi.mock('@tanstack/react-query', () => ({
  useMutation: () => ({
    isPending: false,
    mutate: vi.fn(),
  }),
  useQueryClient: () => ({
    removeQueries: vi.fn(),
  }),
}))

vi.mock('react-hook-form', () => ({
  useForm: () => ({
    formState: { errors: {} },
    handleSubmit: () => (event?: Event) => event,
    register: () => ({}),
    reset: vi.fn(),
    setError: vi.fn(),
  }),
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

vi.mock('../../shared/components/ui/empty-state', () => ({
  EmptyState: ({
    action,
    children,
    title,
  }: { action?: unknown; children?: unknown; title?: string }) => (
    <section>
      <h2>{title}</h2>
      <div>{children}</div>
      {action}
    </section>
  ),
}))

vi.mock('../../shared/components/ui/field', () => ({
  Input: () => <input aria-label="Avatar name" />,
}))

vi.mock('../../shared/components/ui/modal', () => ({
  Modal: ({ children, isOpen }: { children?: unknown; isOpen: boolean }) =>
    isOpen ? <div>{children}</div> : null,
}))

vi.mock('../../shared/components/ui/sidebar-nav', () => ({
  SidebarNav: () => <aside>Sidebar</aside>,
}))

vi.mock('../../shared/components/ui/topbar', () => ({
  Topbar: ({ title }: { title?: string }) => <header>{title}</header>,
}))

vi.mock('react-router-dom', () => ({
  useNavigate: () => navigateMock,
}))

vi.mock('../../queries/useAvatarsQuery', () => ({
  useAvatarsQuery: () => useAvatarsQueryMock(),
  useCreateAvatarMutation: () => useCreateAvatarMutationMock(),
}))

vi.mock('../../queries/useCurrentUserQuery', () => ({
  currentUserQueryKey: ['auth', 'current-user'],
  useCurrentUserQuery: () => useCurrentUserQueryMock(),
}))

describe('CreativeStudioPage', () => {
  beforeEach(() => {
    navigateMock.mockReset()
    useAvatarsQueryMock.mockReset()
    useCreateAvatarMutationMock.mockReset()
    useCurrentUserQueryMock.mockReset()

    useCurrentUserQueryMock.mockReturnValue({
      data: {
        user: {
          avatarUrl: 'https://avatar.example.com/nico.png',
          email: 'nico@example.com',
          id: 'user-1',
          name: 'Nico',
        },
      },
    })

    useCreateAvatarMutationMock.mockReturnValue({
      error: null,
      isError: false,
      isPending: false,
      mutateAsync: vi.fn(),
    })
  })

  it('retries avatar loading from the local error state instead of relying on app-shell flows', async () => {
    const refetchMock = vi.fn()

    useAvatarsQueryMock.mockReturnValue({
      isError: true,
      isLoading: false,
      refetch: refetchMock,
    })

    render(<CreativeStudioPage />)

    fireEvent.click(screen.getByRole('button', { name: /try again/i }))

    expect(refetchMock).toHaveBeenCalledTimes(1)
  })

  it('navigates to the avatar editor when an avatar card is clicked', async () => {
    useAvatarsQueryMock.mockReturnValue({
      data: {
        avatars: [{ id: 'avatar-1', name: 'Studio Hero' }],
      },
      isError: false,
      isLoading: false,
    })

    render(<CreativeStudioPage />)

    fireEvent.click(
      screen.getByRole('button', {
        name: /open studio hero avatar editor/i,
      }),
    )

    expect(navigateMock).toHaveBeenCalledWith(
      '/creative-studio/avatars/avatar-1/avatar',
    )
  })
})
