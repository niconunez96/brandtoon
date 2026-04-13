import { http, HttpResponse } from 'msw'
import { API_BASE_URL } from '../../shared/config/api'

export const handlers = [
  http.get(`${API_BASE_URL}/greeting/:name`, ({ params }) => {
    const name = String(params.name ?? 'world')
    return HttpResponse.json({ message: `Hello, ${name}!` })
  }),
  http.get(`${API_BASE_URL}/auth/users/me`, () => {
    return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
  }),
  http.post(`${API_BASE_URL}/auth/logout`, () => {
    return HttpResponse.json({ message: 'Logged out' })
  }),
  http.get(`${API_BASE_URL}/creative-studio/avatars`, () => {
    return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
  }),
  http.post(`${API_BASE_URL}/creative-studio/avatars`, () => {
    return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
  }),
  http.get(`${API_BASE_URL}/creative-studio/avatar_configs/:avatarId`, () => {
    return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
  }),
  http.put(`${API_BASE_URL}/creative-studio/avatar_configs/:avatarId`, () => {
    return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
  }),
]
