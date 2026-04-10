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
  it('renders the main showcase sections', () => {
    renderApp()

    expect(
      screen.getByText(/Foundational components and palette/i),
    ).toBeInTheDocument()
    expect(screen.getByText('Foundation Rules')).toBeInTheDocument()
    expect(screen.getByText('Base Components')).toBeInTheDocument()
    expect(screen.getByText('Data & Status Foundations')).toBeInTheDocument()
  })

  it('renders both desktop and mobile navigation landmarks', () => {
    renderApp()

    expect(screen.getByLabelText('Primary sidebar')).toBeInTheDocument()
    expect(
      screen.getByLabelText('Mobile bottom navigation'),
    ).toBeInTheDocument()
  })

  it('renders key primitive examples', () => {
    renderApp()

    expect(screen.getAllByText('Primary action')[0]).toBeInTheDocument()
    expect(screen.getByLabelText('Workspace name')).toBeInTheDocument()
    expect(screen.getByText('System status synced')).toBeInTheDocument()
    expect(screen.getByText('No assets yet')).toBeInTheDocument()
    expect(screen.getByText('Completion')).toBeInTheDocument()
    expect(screen.getByText('#vector')).toBeInTheDocument()
    expect(screen.getByText('Project saved')).toBeInTheDocument()
  })
})
