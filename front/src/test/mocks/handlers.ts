import { http, HttpResponse } from 'msw'
import { API_BASE_URL } from '../../shared/config/api'

export const handlers = [
  http.get(`${API_BASE_URL}/greeting/:name`, ({ params }) => {
    const name = String(params.name ?? 'world')
    return HttpResponse.json({ message: `Hello, ${name}!` })
  }),
  http.get(`${API_BASE_URL}/auth/me`, () => {
    return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
  }),
  http.post(`${API_BASE_URL}/auth/logout`, () => {
    return HttpResponse.json({ message: 'Logged out' })
  }),
]
