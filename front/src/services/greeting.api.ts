import { API_BASE_URL } from '../shared/config/api'

type GreetingResponse = {
  message: string
}

export async function fetchGreeting(name: string): Promise<GreetingResponse> {
  const sanitized = encodeURIComponent(name.trim() || 'world')
  const response = await fetch(`${API_BASE_URL}/greeting/${sanitized}`)

  if (!response.ok) {
    throw new Error('Failed to fetch greeting')
  }

  return response.json()
}
