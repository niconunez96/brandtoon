import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import {
  currentUserQueryKey,
  useCurrentUserQuery,
} from '../../queries/useCurrentUserQuery'
import { logoutSession } from '../../services/auth.api'
import { Button } from '../../shared/components/ui/button'

const GENERIC_AVATAR_SRC = `data:image/svg+xml;utf8,${encodeURIComponent(
  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 96 96" fill="none"><rect width="96" height="96" rx="48" fill="#FCE7E7"/><circle cx="48" cy="36" r="18" fill="#FF7A66" opacity="0.9"/><path d="M18 82c3-17 16-26 30-26s27 9 30 26" fill="#FF7A66" opacity="0.9"/></svg>`,
)}`

export function CreativeStudioSidebarFooter() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const currentUserQuery = useCurrentUserQuery()
  const logoutMutation = useMutation({
    mutationFn: logoutSession,
    onSuccess: async () => {
      await queryClient.removeQueries({ queryKey: currentUserQueryKey })
      navigate('/', { replace: true })
    },
  })

  return (
    <div className="rounded-[2rem] border border-[color:var(--color-stroke-soft)] bg-white/90 p-4 shadow-overshoot">
      <div className="flex items-center gap-3">
        <img
          alt="Generic user avatar"
          className="size-12 rounded-full border border-[color:var(--color-stroke-soft)] object-cover"
          src={GENERIC_AVATAR_SRC}
        />

        <div className="min-w-0">
          <p className="text-xs font-black uppercase tracking-section text-ink/55">
            Logged in
          </p>
          <p className="truncate text-sm font-black text-ink">
            {currentUserQuery.data?.user.name ?? 'Brandtoon user'}
          </p>
        </div>
      </div>

      <Button
        className="mt-4 w-full"
        isLoading={logoutMutation.isPending}
        onClick={() => logoutMutation.mutate()}
        variant="ghost"
      >
        Log out
      </Button>
    </div>
  )
}
