import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  fetchAvatarConfig,
  updateAvatarConfig,
  type UpdateAvatarConfigInput,
} from '../services/avatar-config.api'

export const avatarConfigQueryKey = (avatarId: string) =>
  ['creative-studio', 'avatar-config', avatarId] as const

export function useAvatarConfigQuery(avatarId: string) {
  return useQuery({
    enabled: avatarId.length > 0,
    queryFn: () => fetchAvatarConfig(avatarId),
    queryKey: avatarConfigQueryKey(avatarId),
    retry: false,
  })
}

export function useUpdateAvatarConfigMutation(avatarId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: UpdateAvatarConfigInput) =>
      updateAvatarConfig(avatarId, input),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: avatarConfigQueryKey(avatarId),
      })
    },
  })
}
