import { API_BASE_URL } from '../shared/config/api'
import { ApiError } from './auth.api'

export type Avatar = {
  id: string
  name: string
}

export type ListAvatarsResponse = {
  avatars: Avatar[]
}

export type CreateAvatarResponse = {
  avatar: Avatar
}

export async function fetchAvatars(): Promise<ListAvatarsResponse> {
  const response = await fetch(`${API_BASE_URL}/creative-studio/avatars`, {
    credentials: 'include',
  })

  if (!response.ok) {
    throw new ApiError('Failed to fetch avatars', response.status)
  }

  return response.json()
}

export async function createAvatar(
  name: string,
): Promise<CreateAvatarResponse> {
  const response = await fetch(`${API_BASE_URL}/creative-studio/avatars`, {
    body: JSON.stringify({ name }),
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    method: 'POST',
  })

  if (!response.ok) {
    throw new ApiError('Failed to create avatar', response.status)
  }

  return response.json()
}
