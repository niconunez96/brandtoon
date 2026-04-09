type GreetingResponse = {
  message: string
}

const API_BASE_URL = import.meta.env.VITE_API_URL ?? 'http://127.0.0.1:8888'

export async function fetchGreeting(name: string): Promise<GreetingResponse> {
  const sanitized = encodeURIComponent(name.trim() || 'world')
  const response = await fetch(`${API_BASE_URL}/greeting/${sanitized}`)

  if (!response.ok) {
    throw new Error('Failed to fetch greeting')
  }

  return response.json()
}
