import { type QueryClient, useQuery } from '@tanstack/react-query'
import { fetchAvatar } from '../services/avatar.api'

export const avatarQueryKey = (avatarId: string) =>
	['creative-studio', 'avatar', avatarId] as const

export function useAvatarQuery(avatarId: string) {
	return useQuery({
		enabled: avatarId.length > 0,
		queryFn: () => fetchAvatar(avatarId),
		queryKey: avatarQueryKey(avatarId),
		retry: false,
	})
}

export async function invalidateAvatarQuery(
	queryClient: QueryClient,
	avatarId: string,
) {
	await queryClient.invalidateQueries({
		queryKey: avatarQueryKey(avatarId),
	})
}
