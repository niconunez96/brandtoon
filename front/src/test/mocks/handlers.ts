import { http, HttpResponse } from 'msw'

export const handlers = [
  http.get('http://127.0.0.1:8888/greeting/:name', ({ params }) => {
    const name = String(params.name ?? 'world')
    return HttpResponse.json({ message: `Hello, ${name}!` })
  }),
  http.get('http://127.0.0.1:8888/auth/me', () => {
    return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
  }),
  http.post('http://127.0.0.1:8888/auth/logout', () => {
    return HttpResponse.json({ message: 'Logged out' })
  }),
]
