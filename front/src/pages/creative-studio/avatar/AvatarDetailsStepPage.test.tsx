import { act, fireEvent, render, screen, within } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { AvatarDetailsStepPage } from './AvatarDetailsStepPage'

const navigateMock = vi.fn()
const useAvatarConfigQueryMock = vi.fn()
const useAvatarQueryMock = vi.fn()
const useAvatarOptionsQueryMock = vi.fn()
const useUpdateAvatarConfigMutationMock = vi.fn()
const useGenerateAvatarOptionsMutationMock = vi.fn()
const useDeleteAvatarOptionsMutationMock = vi.fn()
const useSelectAvatarOptionMutationMock = vi.fn()

vi.mock('lucide-react', () => ({
  ArrowLeft: () => null,
  Sparkles: () => null,
  WandSparkles: () => null,
}))

vi.mock('react-hook-form', () => ({
  Controller: ({
    name,
    render,
  }: {
    name: string
    render: (input: {
      field: {
        name: string
        onBlur: () => void
        onChange: () => void
        value: string
      }
    }) => unknown
  }) =>
    render({
      field: {
        name,
        onBlur: vi.fn(),
        onChange: vi.fn(),
        value: 'Studio mascot',
      },
    }),
  useForm: () => ({
    clearErrors: vi.fn(),
    control: {},
    formState: { errors: {} },
    handleSubmit: (callback: (values: unknown) => unknown) => () =>
      callback({
        artisticStyle: '2D',
        personality: 'Friendly',
        prompt: 'Studio mascot',
      }),
    register: () => ({}),
    reset: vi.fn(),
    setError: vi.fn(),
    setValue: vi.fn(),
    watch: (field: string) => (field === 'artisticStyle' ? '2D' : 'Friendly'),
  }),
}))

vi.mock('react-router-dom', () => ({
  useNavigate: () => navigateMock,
  useParams: () => ({ avatarId: 'avatar-1' }),
}))

vi.mock('../../../shared/components/ui/button', () => ({
  Button: ({
    children,
    disabled,
    onClick,
  }: {
    children: unknown
    disabled?: boolean
    onClick?: () => void
  }) => (
    <button disabled={disabled} onClick={onClick} type="button">
      {children}
    </button>
  ),
}))

vi.mock('../../../shared/components/ui/card', () => ({
  Card: ({ children }: { children: unknown }) => <div>{children}</div>,
  SectionShell: ({
    actions,
    children,
    title,
  }: {
    actions?: unknown
    children: unknown
    title?: string
  }) => (
    <section>
      <h1>{title}</h1>
      {actions}
      {children}
    </section>
  ),
}))

vi.mock('../../../shared/components/ui/field', () => ({
  PromptField: ({
    onChange,
    title,
  }: {
    onChange?: (event: { target: { value: string } }) => void
    title: string
  }) => (
    <textarea
      aria-label={title}
      onChange={() => onChange?.({ target: { value: 'Updated mascot' } })}
    />
  ),
}))

vi.mock('../../../shared/components/ui/toast', () => ({
  Alert: ({ children, title }: { children: unknown; title: string }) => (
    <div>
      <p>{title}</p>
      <p>{children}</p>
    </div>
  ),
}))

vi.mock('../../../queries/useAvatarConfigQuery', () => ({
  useAvatarConfigQuery: () => useAvatarConfigQueryMock(),
  useDeleteAvatarOptionsMutation: () => useDeleteAvatarOptionsMutationMock(),
  useGenerateAvatarOptionsMutation: () =>
    useGenerateAvatarOptionsMutationMock(),
  useSelectAvatarOptionMutation: () => useSelectAvatarOptionMutationMock(),
  useUpdateAvatarConfigMutation: () => useUpdateAvatarConfigMutationMock(),
}))

vi.mock('../../../queries/useAvatarQuery', () => ({
  useAvatarQuery: () => useAvatarQueryMock(),
}))

vi.mock('../../../queries/useAvatarOptionsQuery', () => ({
  useAvatarOptionsQuery: () => useAvatarOptionsQueryMock(),
}))

function renderAvatarDetailsPage() {
  return render(<AvatarDetailsStepPage />)
}

function buildAvatarConfig() {
  return {
    avatarId: 'avatar-1',
    artisticStyle: '2D' as const,
    personality: 'Friendly' as const,
    prompt: 'Studio mascot',
  }
}

function buildAvatar() {
  return {
    id: 'avatar-1',
    name: 'Studio mascot',
  }
}

function buildAvatarOptions() {
  return {
    avatarId: 'avatar-1',
    avatarOptions: [
      {
        href: 'https://cdn.brandtoon.local/avatars/avatar-1/options/1.png',
        id: 'option-1',
        selected: false,
        status: 'DONE',
      },
      {
        href: 'https://cdn.brandtoon.local/avatars/avatar-1/options/2.png',
        id: 'option-2',
        selected: true,
        status: 'DONE',
      },
      {
        href: 'https://cdn.brandtoon.local/avatars/avatar-1/options/3.png',
        id: 'option-3',
        selected: false,
        status: 'DONE',
      },
    ],
  }
}

function mockLoadedAvatarConfig() {
  useAvatarConfigQueryMock.mockReturnValue({
    data: { avatar_config: buildAvatarConfig() },
    isError: false,
    isLoading: false,
    refetch: vi.fn(),
  })
  useAvatarQueryMock.mockReturnValue({
    data: { avatar: buildAvatar() },
    isError: false,
    isLoading: false,
    refetch: vi.fn(),
  })
  useAvatarOptionsQueryMock.mockReturnValue({
    data: { avatar: buildAvatarOptions() },
    isError: false,
    isLoading: false,
    refetch: vi.fn(),
  })

  useUpdateAvatarConfigMutationMock.mockReturnValue({
    isPending: false,
    mutateAsync: vi.fn(),
  })
  useGenerateAvatarOptionsMutationMock.mockReturnValue({
    isPending: false,
    mutateAsync: vi.fn(),
  })
}

describe('AvatarDetailsStepPage', () => {
  it('selects an option through the unit mutation seam', () => {
    const selectAvatarOptionMutation = {
      isPending: false,
      mutateAsync: vi.fn(),
    }

    mockLoadedAvatarConfig()
    useSelectAvatarOptionMutationMock.mockReturnValue(
      selectAvatarOptionMutation,
    )
    useDeleteAvatarOptionsMutationMock.mockReturnValue({
      isPending: false,
      mutateAsync: vi.fn(),
    })

    renderAvatarDetailsPage()

    fireEvent.click(
      screen.getByRole('button', {
        name: /select avatar option option-2/i,
      }),
    )

    expect(selectAvatarOptionMutation.mutateAsync).toHaveBeenCalledWith(
      'option-2',
    )

    const selectedCard = screen.getByText('Selected', { selector: 'span' })
      .parentElement?.parentElement

    expect(selectedCard).not.toBeNull()
    expect(
      within(selectedCard as HTMLElement).getByRole('button', {
        name: /select avatar option option-2/i,
      }),
    ).toBeInTheDocument()

    const renderedImageSources = screen
      .getAllByRole('img', { name: /avatar option/i })
      .map((image) => image.getAttribute('src'))

    expect(renderedImageSources).toEqual([
      'https://cdn.brandtoon.local/avatars/avatar-1/options/1.png',
      'https://cdn.brandtoon.local/avatars/avatar-1/options/2.png',
      'https://cdn.brandtoon.local/avatars/avatar-1/options/3.png',
    ])
  })

  it('renders avatar options from the dedicated avatar-options query', () => {
    mockLoadedAvatarConfig()
    useAvatarOptionsQueryMock.mockReturnValue({
      data: {
        avatar: {
          avatarOptions: [
            {
              avatarId: 'avatar-1',
              id: 'option-from-avatar',
              href: 'https://cdn.brandtoon.local/avatars/avatar-1/options/avatar-only.png',
              selected: true,
              status: 'DONE',
            },
          ],
        },
      },
      isError: false,
      isLoading: false,
      refetch: vi.fn(),
    })
    useSelectAvatarOptionMutationMock.mockReturnValue({
      isPending: false,
      mutateAsync: vi.fn(),
    })
    useDeleteAvatarOptionsMutationMock.mockReturnValue({
      isPending: false,
      mutateAsync: vi.fn(),
    })

    renderAvatarDetailsPage()

    expect(
      screen.getByRole('button', {
        name: /select avatar option option-from-avatar/i,
      }),
    ).toBeInTheDocument()
    expect(
      screen.queryByRole('button', {
        name: /select avatar option option-2/i,
      }),
    ).not.toBeInTheDocument()
  })

  it('renders the selected avatar image in the large preview box', () => {
    mockLoadedAvatarConfig()
    useSelectAvatarOptionMutationMock.mockReturnValue({
      isPending: false,
      mutateAsync: vi.fn(),
    })
    useDeleteAvatarOptionsMutationMock.mockReturnValue({
      isPending: false,
      mutateAsync: vi.fn(),
    })

    renderAvatarDetailsPage()

    expect(
      screen.getByRole('img', { name: /selected avatar preview/i }),
    ).toHaveAttribute(
      'src',
      'https://cdn.brandtoon.local/avatars/avatar-1/options/2.png',
    )
  })

  it('keeps the preview fallback when no option is selected yet', () => {
    mockLoadedAvatarConfig()
    useAvatarOptionsQueryMock.mockReturnValue({
      data: {
        avatar: {
          avatarOptions: buildAvatarOptions().avatarOptions.map((option) => ({
            ...option,
            selected: false,
          })),
        },
      },
      isError: false,
      isLoading: false,
      refetch: vi.fn(),
    })
    useSelectAvatarOptionMutationMock.mockReturnValue({
      isPending: false,
      mutateAsync: vi.fn(),
    })
    useDeleteAvatarOptionsMutationMock.mockReturnValue({
      isPending: false,
      mutateAsync: vi.fn(),
    })

    renderAvatarDetailsPage()

    expect(
      screen.queryByRole('img', { name: /selected avatar preview/i }),
    ).not.toBeInTheDocument()
    expect(screen.getByText(/active prototype/i)).toBeInTheDocument()
  })

  it('deletes the marked options and clears the local delete counter', async () => {
    const deleteAvatarOptionsMutation = {
      isPending: false,
      mutateAsync: vi.fn().mockResolvedValue(undefined),
    }

    mockLoadedAvatarConfig()
    useSelectAvatarOptionMutationMock.mockReturnValue({
      isPending: false,
      mutateAsync: vi.fn(),
    })
    useDeleteAvatarOptionsMutationMock.mockReturnValue(
      deleteAvatarOptionsMutation,
    )

    renderAvatarDetailsPage()

    fireEvent.click(
      screen.getByLabelText(/mark avatar option option-1 for deletion/i),
    )
    fireEvent.click(
      screen.getByLabelText(/mark avatar option option-3 for deletion/i),
    )
    fireEvent.click(
      screen.getByRole('button', { name: /delete selected \(2\)/i }),
    )

    expect(deleteAvatarOptionsMutation.mutateAsync).toHaveBeenCalledWith([
      'option-1',
      'option-3',
    ])
    expect(
      await screen.findByRole('button', { name: /delete selected \(0\)/i }),
    )
  })

  it('shows draft save success feedback, fades it out, and clears it after editing again', async () => {
    vi.useFakeTimers()

    try {
      const updateAvatarConfigMutation = {
        isPending: false,
        mutateAsync: vi.fn().mockResolvedValue(undefined),
      }

      mockLoadedAvatarConfig()
      useUpdateAvatarConfigMutationMock.mockReturnValue(
        updateAvatarConfigMutation,
      )
      useSelectAvatarOptionMutationMock.mockReturnValue({
        isPending: false,
        mutateAsync: vi.fn(),
      })
      useDeleteAvatarOptionsMutationMock.mockReturnValue({
        isPending: false,
        mutateAsync: vi.fn(),
      })

      renderAvatarDetailsPage()

      await act(async () => {
        fireEvent.click(screen.getByRole('button', { name: /save as draft/i }))
      })

      expect(
        screen.getByText(/your avatar draft was saved\./i),
      ).toBeInTheDocument()

      act(() => {
        vi.advanceTimersByTime(2600)
      })

      expect(
        screen.queryByText(/your avatar draft was saved\./i),
      ).not.toBeInTheDocument()

      fireEvent.click(screen.getByRole('button', { name: '3D' }))

      expect(
        screen.queryByText(/your avatar draft was saved\./i),
      ).not.toBeInTheDocument()
    } finally {
      vi.useRealTimers()
    }
  })

  it('disables delete actions while selection is pending', () => {
    mockLoadedAvatarConfig()
    useSelectAvatarOptionMutationMock.mockReturnValue({
      isPending: true,
      mutateAsync: vi.fn(),
    })
    useDeleteAvatarOptionsMutationMock.mockReturnValue({
      isPending: false,
      mutateAsync: vi.fn(),
    })

    renderAvatarDetailsPage()

    fireEvent.click(
      screen.getByLabelText(/mark avatar option option-1 for deletion/i),
    )

    expect(
      screen.getByRole('button', { name: /delete selected \(1\)/i }),
    ).toBeDisabled()
    expect(
      screen.getByLabelText(/mark avatar option option-2 for deletion/i),
    ).toBeDisabled()
  })

  it('disables selection while delete is pending', () => {
    mockLoadedAvatarConfig()
    useSelectAvatarOptionMutationMock.mockReturnValue({
      isPending: false,
      mutateAsync: vi.fn(),
    })
    useDeleteAvatarOptionsMutationMock.mockReturnValue({
      isPending: true,
      mutateAsync: vi.fn(),
    })

    renderAvatarDetailsPage()

    expect(
      screen.getByRole('button', { name: /select avatar option option-1/i }),
    ).toBeDisabled()
    expect(
      screen.getByLabelText(/mark avatar option option-1 for deletion/i),
    ).toBeDisabled()
  })
})
