import { act, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { AvatarGenerationEventsProvider } from './AvatarGenerationEventsProvider'

const invalidateAvatarOptionsQueryMock = vi.fn()
const addEventListenerMock = vi.fn()
const removeEventListenerMock = vi.fn()
const closeMock = vi.fn()
const useLocationMock = vi.fn()

vi.mock('@tanstack/react-query', () => ({
  useQueryClient: () => ({ queryClient: true }),
}))

vi.mock('react-router-dom', async () => {
  const actual =
    await vi.importActual<typeof import('react-router-dom')>('react-router-dom')

  return {
    ...actual,
    useLocation: () => useLocationMock(),
  }
})

vi.mock('../../queries/useAvatarOptionsQuery', () => ({
  invalidateAvatarOptionsQuery: (...args: unknown[]) =>
    invalidateAvatarOptionsQueryMock(...args),
}))

vi.mock('../../shared/components/ui/toast', () => ({
  Toast: ({ children, title }: { children: unknown; title: string }) => (
    <div>
      <p>{title}</p>
      <p>{children}</p>
    </div>
  ),
}))

describe('AvatarGenerationEventsProvider', () => {
  afterEach(() => {
    invalidateAvatarOptionsQueryMock.mockReset()
    addEventListenerMock.mockReset()
    removeEventListenerMock.mockReset()
    closeMock.mockReset()
    useLocationMock.mockReset()
    vi.unstubAllGlobals()
  })

  it('invalidates the avatar-options query when generation completes for the active avatar', () => {
    let completedHandler: ((event: MessageEvent<string>) => void) | undefined

    class EventSourceMock {
      addEventListener(
        type: string,
        handler: (event: MessageEvent<string>) => void,
      ) {
        addEventListenerMock(type)
        if (type === 'avatar-generation.completed') {
          completedHandler = handler
        }
      }

      removeEventListener(type: string) {
        removeEventListenerMock(type)
      }

      close() {
        closeMock()
      }
    }

    vi.stubGlobal('EventSource', EventSourceMock)
    useLocationMock.mockReturnValue({
      pathname: '/creative-studio/avatars/avatar-v7/avatar',
    })

    render(
      <AvatarGenerationEventsProvider>
        <div>child</div>
      </AvatarGenerationEventsProvider>,
    )

    act(() => {
      completedHandler?.({
        data: JSON.stringify({
          avatarId: 'avatar-v7',
          avatarName: 'Studio Hero',
          outcome: 'SUCCESS',
          userId: 'user-v7',
        }),
      } as MessageEvent<string>)
    })

    expect(invalidateAvatarOptionsQueryMock).toHaveBeenCalledWith(
      expect.anything(),
      'avatar-v7',
    )
  })

  it('shows success messaging only for successful completion events', () => {
    let completedHandler: ((event: MessageEvent<string>) => void) | undefined

    class EventSourceMock {
      addEventListener(
        type: string,
        handler: (event: MessageEvent<string>) => void,
      ) {
        if (type === 'avatar-generation.completed') {
          completedHandler = handler
        }
      }

      removeEventListener() {}

      close() {}
    }

    vi.stubGlobal('EventSource', EventSourceMock)
    useLocationMock.mockReturnValue({ pathname: '/creative-studio' })

    render(
      <AvatarGenerationEventsProvider>
        <div>child</div>
      </AvatarGenerationEventsProvider>,
    )

    act(() => {
      completedHandler?.({
        data: JSON.stringify({
          avatarId: 'avatar-v7',
          avatarName: 'Studio Hero',
          outcome: 'SUCCESS',
          userId: 'user-v7',
        }),
      } as MessageEvent<string>)
    })

    expect(
      screen.getByText('Avatar generation succeeded for Studio Hero'),
    ).toBeInTheDocument()
    expect(screen.getByText('Your avatar options are ready.')).toBeInTheDocument()
    expect(screen.queryByText(/failed/i)).not.toBeInTheDocument()
  })

  it('shows failure messaging for failed completion events', () => {
    let completedHandler: ((event: MessageEvent<string>) => void) | undefined

    class EventSourceMock {
      addEventListener(
        type: string,
        handler: (event: MessageEvent<string>) => void,
      ) {
        if (type === 'avatar-generation.completed') {
          completedHandler = handler
        }
      }

      removeEventListener() {}

      close() {}
    }

    vi.stubGlobal('EventSource', EventSourceMock)
    useLocationMock.mockReturnValue({ pathname: '/creative-studio' })

    render(
      <AvatarGenerationEventsProvider>
        <div>child</div>
      </AvatarGenerationEventsProvider>,
    )

    act(() => {
      completedHandler?.({
        data: JSON.stringify({
          avatarId: 'avatar-v7',
          avatarName: 'Studio Hero',
          outcome: 'FAILURE',
          userId: 'user-v7',
        }),
      } as MessageEvent<string>)
    })

    expect(
      screen.getByText('Avatar generation failed for Studio Hero'),
    ).toBeInTheDocument()
    expect(
      screen.getByText(
        'We could not generate avatar options this time. Please try again.',
      ),
    ).toBeInTheDocument()
    expect(screen.queryByText(/succeeded/i)).not.toBeInTheDocument()
  })
})
