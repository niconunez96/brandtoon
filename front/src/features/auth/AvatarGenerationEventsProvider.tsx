import { useQueryClient } from '@tanstack/react-query'
import { Sparkles } from 'lucide-react'
import { type PropsWithChildren, useEffect, useMemo, useState } from 'react'
import { matchPath, useLocation } from 'react-router-dom'
import { invalidateAvatarConfigQuery } from '../../queries/useAvatarConfigQuery'
import type { AvatarGenerationCompletedEvent } from '../../services/avatar-config.api'
import { Toast } from '../../shared/components/ui/toast'
import { API_BASE_URL } from '../../shared/config/api'

type ToastVisibility = 'hidden' | 'entering' | 'visible' | 'exiting'

export function AvatarGenerationEventsProvider({
  children,
}: PropsWithChildren) {
  const location = useLocation()
  const queryClient = useQueryClient()
  const [toastMessage, setToastMessage] = useState<string | null>(null)
  const [toastVisibility, setToastVisibility] =
    useState<ToastVisibility>('hidden')

  const activeAvatarId = useMemo(() => {
    const match = matchPath(
      '/creative-studio/avatars/:avatarId/avatar',
      location.pathname,
    )

    return match?.params.avatarId ?? null
  }, [location.pathname])

  useEffect(() => {
    const eventSource = new EventSource(`${API_BASE_URL}/events`, {
      withCredentials: true,
    })

    const handleCompleted = (event: MessageEvent<string>) => {
      const payload = JSON.parse(event.data) as AvatarGenerationCompletedEvent
      setToastMessage(`The avatar ${payload.avatarName} was generated`)
      setToastVisibility('entering')

      if (activeAvatarId === payload.avatarId) {
        void invalidateAvatarConfigQuery(queryClient, payload.avatarId)
      }
    }

    eventSource.addEventListener('avatar-generation.completed', handleCompleted)

    return () => {
      eventSource.removeEventListener(
        'avatar-generation.completed',
        handleCompleted,
      )
      eventSource.close()
    }
  }, [activeAvatarId, queryClient])

  useEffect(() => {
    if (!toastMessage || toastVisibility === 'hidden') {
      return
    }

    if (toastVisibility === 'entering') {
      const enterTimer = window.setTimeout(
        () => setToastVisibility('visible'),
        16,
      )
      return () => window.clearTimeout(enterTimer)
    }

    if (toastVisibility === 'visible') {
      const visibleTimer = window.setTimeout(
        () => setToastVisibility('exiting'),
        2400,
      )
      return () => window.clearTimeout(visibleTimer)
    }

    const exitTimer = window.setTimeout(() => {
      setToastVisibility('hidden')
      setToastMessage(null)
    }, 250)

    return () => window.clearTimeout(exitTimer)
  }, [toastMessage, toastVisibility])

  return (
    <>
      {children}
      {toastMessage && toastVisibility !== 'hidden' ? (
        <div className="pointer-events-none fixed right-4 top-4 z-50 sm:right-6 sm:top-6">
          <div
            className={`pointer-events-auto transition-all duration-250 ease-out ${
              toastVisibility === 'visible'
                ? 'translate-x-0 opacity-100'
                : '-translate-x-6 opacity-0'
            }`}
          >
            <Toast
              icon={<Sparkles className="size-4" />}
              onDismiss={() => setToastVisibility('exiting')}
              title={toastMessage}
            >
              Your avatar options are ready.
            </Toast>
          </div>
        </div>
      ) : null}
    </>
  )
}
