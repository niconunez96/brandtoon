import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { fireEvent, render, screen } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { App } from './App'
import { useGreetingStore } from './store/greeting.store'
import { server } from './test/mocks/server'

const renderApp = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  })

  return render(
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>,
  )
}

describe('App', () => {
  afterEach(() => {
    useGreetingStore.setState({ name: 'world' })
  })

  it('renders greeting from API', async () => {
    renderApp()

    expect(screen.getByText('Loading greeting...')).toBeInTheDocument()
    expect(await screen.findByText('Hello, world!')).toBeInTheDocument()
  })

  it('shows error when API request fails', async () => {
    server.use(
      http.get('http://127.0.0.1:8888/greeting/:name', () =>
        HttpResponse.json({}, { status: 500 }),
      ),
    )

    renderApp()

    expect(
      await screen.findByText('Could not fetch greeting.'),
    ).toBeInTheDocument()
  })

  it('updates greeting when user changes the name', async () => {
    renderApp()

    const input = await screen.findByLabelText('Name')
    fireEvent.change(input, { target: { value: 'nico' } })

    expect(await screen.findByText('Hello, nico!')).toBeInTheDocument()
  })
})
