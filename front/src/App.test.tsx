import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen } from '@testing-library/react'
import { App } from './App'

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
  it('renders the landing page sections and hero content', () => {
    renderApp()

    expect(
      screen.getByRole('heading', { name: /your brand, animated\./i }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('link', { name: 'How it Works' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: /create your avatar/i }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: /watch showcase/i }),
    ).toBeInTheDocument()
  })

  it('renders pricing and final call to action content', () => {
    renderApp()

    expect(
      screen.getByRole('heading', { name: /no strings attached\./i }),
    ).toBeInTheDocument()
    expect(
      screen.getByText(/includes 50\+ custom variations/i),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: /purchase credits/i }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('heading', { name: /ready to toon your brand\?/i }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: /launch generator/i }),
    ).toBeInTheDocument()
  })

  it('removes the old component catalog content', () => {
    renderApp()

    expect(screen.queryByText('Base Components')).not.toBeInTheDocument()
    expect(
      screen.queryByText('Data & Status Foundations'),
    ).not.toBeInTheDocument()
    expect(screen.queryByLabelText('Primary sidebar')).not.toBeInTheDocument()
  })
})
