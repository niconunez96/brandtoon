import { API_BASE_URL } from '../shared/config/api'
import { ApiError } from './auth.api'

export type AvatarOptionStatus = 'PENDING' | 'DONE' | 'FAILED'

export type AvatarOption = {
  id: string
  status: AvatarOptionStatus
  selected: boolean
  href?: string
}

export type AvatarOptionsResponse = {
  avatar: {
    avatarId: string
    avatarOptions: AvatarOption[]
  } | null
}

export async function fetchAvatarOptions(
  avatarId: string,
): Promise<AvatarOptionsResponse> {
  const response = await fetch(
    `${API_BASE_URL}/creative-studio/avatars/${avatarId}/options`,
    { credentials: 'include' },
  )

  if (!response.ok) {
    throw new ApiError('Failed to fetch avatar options', response.status)
  }

  return response.json()
}

export async function generateAvatarOptions(avatarId: string): Promise<void> {
  const response = await fetch(
    `${API_BASE_URL}/creative-studio/avatars/${avatarId}/options/generate`,
    { credentials: 'include', method: 'POST' },
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
    `${API_BASE_URL}/creative-studio/avatars/${avatarId}/options/select`,
    {
      body: JSON.stringify({ id }),
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
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
    `${API_BASE_URL}/creative-studio/avatars/${avatarId}/options`,
    {
      body: JSON.stringify({ ids }),
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      method: 'DELETE',
    },
  )

  if (!response.ok) {
    throw new ApiError('Failed to delete avatar options', response.status)
  }
}
