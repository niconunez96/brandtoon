import { http, HttpResponse } from 'msw'

export const handlers = [
  http.get('http://127.0.0.1:8888/greeting/:name', ({ params }) => {
    const name = String(params.name ?? 'world')
    return HttpResponse.json({ message: `Hello, ${name}!` })
  }),
]
