export const API_BASE_URL =
  import.meta.env.VITE_API_URL ?? 'http://127.0.0.1:8888'

export function sanitizeNextPath(nextPath: string | null | undefined) {
  const trimmed = nextPath?.trim()

  if (!trimmed || !trimmed.startsWith('/') || trimmed.startsWith('//')) {
    return '/creative-studio'
  }

  return trimmed
}
