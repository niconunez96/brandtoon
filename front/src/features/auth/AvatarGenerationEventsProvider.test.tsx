import { act, render } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { AvatarGenerationEventsProvider } from './AvatarGenerationEventsProvider'

const invalidateAvatarQueryMock = vi.fn()
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

vi.mock('../../queries/useAvatarQuery', () => ({
  invalidateAvatarQuery: (...args: unknown[]) =>
    invalidateAvatarQueryMock(...args),
}))

vi.mock('../../shared/components/ui/toast', () => ({
  Toast: ({ children }: { children: unknown }) => <div>{children}</div>,
}))

describe('AvatarGenerationEventsProvider', () => {
  afterEach(() => {
    invalidateAvatarQueryMock.mockReset()
    addEventListenerMock.mockReset()
    removeEventListenerMock.mockReset()
    closeMock.mockReset()
    useLocationMock.mockReset()
    vi.unstubAllGlobals()
  })

  it('invalidates the avatar query when generation completes for the active avatar', () => {
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
          userId: 'user-v7',
        }),
      } as MessageEvent<string>)
    })

    expect(invalidateAvatarQueryMock).toHaveBeenCalledWith(
      expect.anything(),
      'avatar-v7',
    )
  })
})
