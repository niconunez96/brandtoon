import {
  type QueryClient,
  useMutation,
  useQuery,
  useQueryClient,
} from '@tanstack/react-query'
import {
  deleteAvatarOptions,
  selectAvatarOption,
  type UpdateAvatarConfigInput,
  fetchAvatarConfig,
  generateAvatarOptions,
  updateAvatarConfig,
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

export async function invalidateAvatarConfigQuery(
  queryClient: QueryClient,
  avatarId: string,
) {
  await queryClient.invalidateQueries({
    queryKey: avatarConfigQueryKey(avatarId),
  })
}

export function useGenerateAvatarOptionsMutation(avatarId: string) {
  return useMutation({
    mutationFn: () => generateAvatarOptions(avatarId),
  })
}

export function useSelectAvatarOptionMutation(avatarId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => selectAvatarOption(avatarId, id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: avatarConfigQueryKey(avatarId),
      })
    },
  })
}

export function useDeleteAvatarOptionsMutation(avatarId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (ids: string[]) => deleteAvatarOptions(avatarId, ids),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: avatarConfigQueryKey(avatarId),
      })
    },
  })
}
