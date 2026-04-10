import { useQuery } from '@tanstack/react-query'
import { fetchCurrentUser, isUnauthorizedError } from '../services/auth.api'

export const currentUserQueryKey = ['auth', 'current-user'] as const

export function useCurrentUserQuery() {
  return useQuery({
    queryKey: currentUserQueryKey,
    queryFn: fetchCurrentUser,
    retry: (failureCount, error) => {
      if (isUnauthorizedError(error)) {
        return false
      }

      return failureCount < 1
    },
    staleTime: 60_000,
  })
}
