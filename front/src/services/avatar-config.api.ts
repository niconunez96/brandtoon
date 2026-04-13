import { API_BASE_URL } from '../shared/config/api'
import { ApiError } from './auth.api'

export type ArtisticStyle = '2D' | '3D'

export type AvatarConfig = {
	avatarId: string
	artisticStyle: ArtisticStyle
	prompt: string
}

export type AvatarConfigResponse = {
	avatar_config: AvatarConfig | null
}

export type UpdateAvatarConfigInput = {
	artisticStyle: ArtisticStyle
	prompt: string
}

export async function fetchAvatarConfig(
	avatarId: string,
): Promise<AvatarConfigResponse> {
	const response = await fetch(
		`${API_BASE_URL}/creative-studio/avatar_configs/${avatarId}`,
		{
			credentials: 'include',
		},
	)

	if (!response.ok) {
		throw new ApiError('Failed to fetch avatar config', response.status)
	}

	return response.json()
}

export async function updateAvatarConfig(
	avatarId: string,
	input: UpdateAvatarConfigInput,
): Promise<AvatarConfigResponse> {
	const response = await fetch(
		`${API_BASE_URL}/creative-studio/avatar_configs/${avatarId}`,
		{
			body: JSON.stringify(input),
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json',
			},
			method: 'PUT',
		},
	)

	if (!response.ok) {
		throw new ApiError('Failed to update avatar config', response.status)
	}

	return response.json()
}
