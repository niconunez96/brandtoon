import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import {
  currentUserQueryKey,
  useCurrentUserQuery,
} from '../../queries/useCurrentUserQuery'
import { logoutSession } from '../../services/auth.api'
import { Button } from '../../shared/components/ui/button'

export function CreativeStudioPage() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  useCurrentUserQuery()
  const logoutMutation = useMutation({
    mutationFn: logoutSession,
    onSuccess: async () => {
      await queryClient.removeQueries({ queryKey: currentUserQueryKey })
      navigate('/', { replace: true })
    },
  })

  return (
    <div className="foundation-page min-h-screen">
      <header className="border-b border-[color:var(--color-stroke-soft)] bg-page/92 px-4 py-5 backdrop-blur-xl sm:px-6 lg:px-8">
        <nav
          aria-label="Creative studio navigation"
          className="mx-auto flex w-full max-w-6xl justify-end"
        >
          <Button
            isLoading={logoutMutation.isPending}
            onClick={() => logoutMutation.mutate()}
            variant="ghost"
          >
            Log out
          </Button>
        </nav>
      </header>

      <main className="mx-auto min-h-[28rem] max-w-6xl px-4 pb-12 pt-6 sm:px-6 lg:px-8" />
    </div>
  )
}
