import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createAvatar, fetchAvatars } from '../services/avatar.api'

export const avatarsQueryKey = ['creative-studio', 'avatars'] as const

export function useAvatarsQuery() {
  return useQuery({
    queryKey: avatarsQueryKey,
    queryFn: fetchAvatars,
    retry: false,
  })
}

export function useCreateAvatarMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: createAvatar,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: avatarsQueryKey })
    },
  })
}
