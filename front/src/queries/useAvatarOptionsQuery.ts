import { type QueryClient, useQuery } from '@tanstack/react-query'
import { fetchAvatarOptions } from '../services/avatar-option.api'

export const avatarOptionsQueryKey = (avatarId: string) =>
  ['creative-studio', 'avatar-options', avatarId] as const

export function useAvatarOptionsQuery(avatarId: string) {
  return useQuery({
    enabled: avatarId.length > 0,
    queryFn: () => fetchAvatarOptions(avatarId),
    queryKey: avatarOptionsQueryKey(avatarId),
    retry: false,
  })
}

export async function invalidateAvatarOptionsQuery(
  queryClient: QueryClient,
  avatarId: string,
) {
  await queryClient.invalidateQueries({
    queryKey: avatarOptionsQueryKey(avatarId),
  })
}
