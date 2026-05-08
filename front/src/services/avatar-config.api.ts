import { API_BASE_URL } from '../shared/config/api'
import { ApiError } from './auth.api'

export type ArtisticStyle = '2D' | '3D'
export type Personality = 'Friendly' | 'Bold' | 'Playful'

export type AvatarConfig = {
  avatarId: string
  artisticStyle: ArtisticStyle
  personality: Personality
  prompt: string
}

export type AvatarGenerationCompletedEvent = {
  avatarId: string
  avatarName: string
  userId: string
}

export type AvatarConfigResponse = {
  avatar_config: AvatarConfig | null
}

export type UpdateAvatarConfigInput = {
  artisticStyle: ArtisticStyle
  personality: Personality
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

export async function generateAvatarOptions(avatarId: string): Promise<void> {
  const response = await fetch(
    `${API_BASE_URL}/creative-studio/avatar_configs/${avatarId}/generate`,
    {
      credentials: 'include',
      method: 'POST',
    },
  )

  if (!response.ok) {
    throw new ApiError('Failed to generate avatar options', response.status)
  }
}

export async function selectAvatarOption(
  avatarId: string,
  id: string,
): Promise<void> {
  const response = await fetch(
    `${API_BASE_URL}/creative-studio/avatar_configs/${avatarId}/options/select`,
    {
      body: JSON.stringify({ id }),
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      method: 'POST',
    },
  )

  if (!response.ok) {
    throw new ApiError('Failed to select avatar option', response.status)
  }
}

export async function deleteAvatarOptions(
  avatarId: string,
  ids: string[],
): Promise<void> {
  const response = await fetch(
    `${API_BASE_URL}/creative-studio/avatar_configs/${avatarId}/options`,
    {
      body: JSON.stringify({ ids }),
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      method: 'DELETE',
    },
  )

  if (!response.ok) {
    throw new ApiError('Failed to delete avatar options', response.status)
  }
}
