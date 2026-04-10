import { API_BASE_URL, sanitizeNextPath } from '../shared/config/api'

export type AuthUser = {
  avatarUrl: string
  email: string
  id: string
  name: string
}

export type CurrentUserResponse = {
  user: AuthUser
}

export class ApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

export function buildGoogleLoginUrl(nextPath?: string | null) {
  const redirectTo = sanitizeNextPath(nextPath)
  const searchParams = new URLSearchParams({ redirectTo })

  return `${API_BASE_URL}/auth/google/login?${searchParams.toString()}`
}

export async function fetchCurrentUser(): Promise<CurrentUserResponse> {
  const response = await fetch(`${API_BASE_URL}/auth/me`, {
    credentials: 'include',
  })

  if (!response.ok) {
    throw new ApiError('Failed to fetch current user', response.status)
  }

  return response.json()
}

export async function logoutSession() {
  const response = await fetch(`${API_BASE_URL}/auth/logout`, {
    credentials: 'include',
    method: 'POST',
  })

  if (!response.ok) {
    throw new ApiError('Failed to log out', response.status)
  }

  return response.json() as Promise<{ message: string }>
}

export function isUnauthorizedError(error: unknown) {
  return error instanceof ApiError && error.status === 401
}
