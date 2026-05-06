import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { AvatarDetailsStepPage } from './AvatarDetailsStepPage'

const useAvatarConfigQueryMock = vi.fn()
const useUpdateAvatarConfigMutationMock = vi.fn()
const useGenerateAvatarOptionsMutationMock = vi.fn()
const useDeleteAvatarOptionsMutationMock = vi.fn()
const useSelectAvatarOptionMutationMock = vi.fn()

vi.mock('lucide-react', () => ({
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
    handleSubmit: (callback: (values: unknown) => unknown) =>
      () => callback({
        artisticStyle: '2D',
        personality: 'Friendly',
        prompt: 'Studio mascot',
      }),
    register: () => ({}),
    reset: vi.fn(),
    setError: vi.fn(),
    setValue: vi.fn(),
    watch: (field: string) =>
      field === 'artisticStyle' ? '2D' : 'Friendly',
  }),
}))

vi.mock('react-router-dom', () => ({
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
  PromptField: ({ title }: { title: string }) => <textarea aria-label={title} />,
}))

vi.mock('../../../queries/useAvatarConfigQuery', () => ({
  useAvatarConfigQuery: () => useAvatarConfigQueryMock(),
  useDeleteAvatarOptionsMutation: () => useDeleteAvatarOptionsMutationMock(),
  useGenerateAvatarOptionsMutation: () =>
    useGenerateAvatarOptionsMutationMock(),
  useSelectAvatarOptionMutation: () => useSelectAvatarOptionMutationMock(),
  useUpdateAvatarConfigMutation: () => useUpdateAvatarConfigMutationMock(),
}))

function renderAvatarDetailsPage() {
  return render(<AvatarDetailsStepPage />)
}

function buildAvatarConfig() {
  return {
    avatarId: 'avatar-1',
    artisticStyle: '2D' as const,
    avatarOptions: [
      {
        id: 'option-1',
        href: 'https://cdn.brandtoon.local/avatars/avatar-1/options/1.png',
        selected: false,
      },
      {
        id: 'option-2',
        href: 'https://cdn.brandtoon.local/avatars/avatar-1/options/2.png',
        selected: true,
      },
      {
        id: 'option-3',
        href: 'https://cdn.brandtoon.local/avatars/avatar-1/options/3.png',
        selected: false,
      },
    ],
    personality: 'Friendly' as const,
    prompt: 'Studio mascot',
  }
}

function mockLoadedAvatarConfig() {
  useAvatarConfigQueryMock.mockReturnValue({
    data: { avatar_config: buildAvatarConfig() },
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
    useSelectAvatarOptionMutationMock.mockReturnValue(selectAvatarOptionMutation)
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
    expect(screen.getByText('Selected', { selector: 'span' })).toBeInTheDocument()
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
})
