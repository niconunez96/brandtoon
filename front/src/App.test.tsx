import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { http, HttpResponse, delay } from 'msw'
import { MemoryRouter, useLocation } from 'react-router-dom'
import { describe, expect, it, vi } from 'vitest'
import { App } from './App'
import { API_BASE_URL } from './shared/config/api'
import { navigateToExternalUrl } from './shared/lib/browser'
import { server } from './test/mocks/server'

vi.mock('./shared/lib/browser', () => ({
  navigateToExternalUrl: vi.fn(),
}))

function LocationProbe() {
  const location = useLocation()

  return (
    <div data-testid="location-display">{`${location.pathname}${location.search}`}</div>
  )
}

const renderApp = (initialEntry = '/') => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  })

  return render(
    <MemoryRouter initialEntries={[initialEntry]}>
      <QueryClientProvider client={queryClient}>
        <App />
        <LocationProbe />
      </QueryClientProvider>
    </MemoryRouter>,
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

  it('navigates landing start-creating CTA clicks to /login?next=/creative-studio', async () => {
    const user = userEvent.setup()

    renderApp()

    await user.click(
      screen.getByRole('button', { name: /create your avatar/i }),
    )

    expect(screen.getByTestId('location-display')).toHaveTextContent(
      '/login?next=%2Fcreative-studio',
    )
    expect(
      screen.getByRole('heading', { name: /log in with google/i }),
    ).toBeInTheDocument()
  })

  it('renders the login page when visiting /login', () => {
    renderApp('/login?next=%2Fcreative-studio')

    expect(
      screen.getByRole('heading', { name: /log in with google/i }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: /continue with google/i }),
    ).toBeInTheDocument()
  })

  it('redirects protected routes to the login page when unauthenticated', async () => {
    renderApp('/creative-studio')

    expect(
      await screen.findByRole('heading', { name: /log in with google/i }),
    ).toBeInTheDocument()
    expect(screen.getByTestId('location-display')).toHaveTextContent(
      '/login?next=%2Fcreative-studio',
    )
  })

  it('navigates to the backend google auth login endpoint from the login CTA', async () => {
    const user = userEvent.setup()
    const navigateToExternalUrlMock = vi.mocked(navigateToExternalUrl)
    navigateToExternalUrlMock.mockReset()

    renderApp('/login?next=%2Fcreative-studio')

    await user.click(
      screen.getByRole('button', { name: /continue with google/i }),
    )

    expect(navigateToExternalUrlMock).toHaveBeenCalledWith(
      `${API_BASE_URL}/auth/google/login?redirectTo=%2Fcreative-studio`,
    )
  })

  it('renders the creative studio shell when the user is authenticated', async () => {
    server.use(
      http.get(`${API_BASE_URL}/auth/users/me`, () => {
        return HttpResponse.json({
          user: {
            avatarUrl: 'https://avatar.example.com/nico.png',
            email: 'nico@example.com',
            id: 'user-v7',
            name: 'Nico',
          },
        })
      }),
      http.get(`${API_BASE_URL}/creative-studio/avatars`, async () => {
        await delay(30)
        return HttpResponse.json({ avatars: [] })
      }),
    )

    renderApp('/creative-studio')

    expect(await screen.findByText(/loading avatars/i)).toBeInTheDocument()

    const navigation = await screen.findByRole('navigation', {
      name: /creative studio navigation/i,
    })

    expect(
      within(navigation).getByRole('button', { name: /log out/i }),
    ).toBeInTheDocument()
    expect(within(navigation).getAllByRole('button')).toHaveLength(1)
    expect(within(navigation).queryByText(/brandtoon/i)).not.toBeInTheDocument()
    expect(
      await screen.findByRole('button', { name: /create avatar now/i }),
    ).toBeInTheDocument()
  })

  it('renders existing avatar cards and create CTA when avatars already exist', async () => {
    server.use(
      http.get(`${API_BASE_URL}/auth/users/me`, () => {
        return HttpResponse.json({
          user: {
            avatarUrl: 'https://avatar.example.com/nico.png',
            email: 'nico@example.com',
            id: 'user-v7',
            name: 'Nico',
          },
        })
      }),
      http.get(`${API_BASE_URL}/creative-studio/avatars`, () => {
        return HttpResponse.json({
          avatars: [
            { id: 'avatar-1', name: 'Studio Hero' },
            { id: 'avatar-2', name: 'Brand Host' },
          ],
        })
      }),
    )

    renderApp('/creative-studio')

    expect(
      await screen.findByRole('heading', { name: /^studio hero$/i }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('heading', { name: /^brand host$/i }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: /^create avatar$/i }),
    ).toBeInTheDocument()
  })

  it('renders a retry state when avatars fail to load', async () => {
    server.use(
      http.get(`${API_BASE_URL}/auth/users/me`, () => {
        return HttpResponse.json({
          user: {
            avatarUrl: 'https://avatar.example.com/nico.png',
            email: 'nico@example.com',
            id: 'user-v7',
            name: 'Nico',
          },
        })
      }),
      http.get(`${API_BASE_URL}/creative-studio/avatars`, () => {
        return HttpResponse.json({ message: 'boom' }, { status: 500 })
      }),
    )

    renderApp('/creative-studio')

    expect(
      await screen.findByText(/we could not load your avatars/i),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: /try again/i }),
    ).toBeInTheDocument()
  })

  it('shows create avatar errors without leaving the page', async () => {
    const user = userEvent.setup()

    server.use(
      http.get(`${API_BASE_URL}/auth/users/me`, () => {
        return HttpResponse.json({
          user: {
            avatarUrl: 'https://avatar.example.com/nico.png',
            email: 'nico@example.com',
            id: 'user-v7',
            name: 'Nico',
          },
        })
      }),
      http.get(`${API_BASE_URL}/creative-studio/avatars`, () => {
        return HttpResponse.json({ avatars: [] })
      }),
      http.post(`${API_BASE_URL}/creative-studio/avatars`, () => {
        return HttpResponse.json({ message: 'invalid' }, { status: 422 })
      }),
    )

    renderApp('/creative-studio')

    await user.click(
      await screen.findByRole('button', { name: /create avatar now/i }),
    )
    await user.type(screen.getByLabelText(/avatar name/i), 'Studio Hero')
    await user.click(screen.getByRole('button', { name: /save avatar/i }))

    expect(
      await screen.findByText(
        /please enter a valid avatar name before saving/i,
      ),
    ).toBeInTheDocument()
  })

  it('refreshes the avatar list after creating a new avatar successfully', async () => {
    const user = userEvent.setup()
    let avatars = [] as Array<{ id: string; name: string }>

    server.use(
      http.get(`${API_BASE_URL}/auth/users/me`, () => {
        return HttpResponse.json({
          user: {
            avatarUrl: 'https://avatar.example.com/nico.png',
            email: 'nico@example.com',
            id: 'user-v7',
            name: 'Nico',
          },
        })
      }),
      http.get(`${API_BASE_URL}/creative-studio/avatars`, () => {
        return HttpResponse.json({ avatars })
      }),
      http.post(
        `${API_BASE_URL}/creative-studio/avatars`,
        async ({ request }) => {
          const body = (await request.json()) as { name: string }
          const nextAvatar = { id: 'avatar-1', name: body.name.trim() }
          avatars = [nextAvatar]
          return HttpResponse.json({ avatar: nextAvatar }, { status: 201 })
        },
      ),
    )

    renderApp('/creative-studio')

    await user.click(
      await screen.findByRole('button', { name: /create avatar now/i }),
    )
    await user.type(screen.getByLabelText(/avatar name/i), 'Studio Hero')
    await user.click(screen.getByRole('button', { name: /save avatar/i }))

    expect(
      await screen.findByRole('heading', { name: /^studio hero$/i }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: /^create avatar$/i }),
    ).toBeInTheDocument()
  })

  it('logs the user out from Creative Studio and redirects back to landing', async () => {
    const user = userEvent.setup()

    server.use(
      http.get(`${API_BASE_URL}/auth/users/me`, () => {
        return HttpResponse.json({
          user: {
            avatarUrl: 'https://avatar.example.com/nico.png',
            email: 'nico@example.com',
            id: 'user-v7',
            name: 'Nico',
          },
        })
      }),
      http.post(`${API_BASE_URL}/auth/logout`, () => {
        return HttpResponse.json({ message: 'Logged out' })
      }),
    )

    renderApp('/creative-studio')

    await user.click(await screen.findByRole('button', { name: /log out/i }))

    expect(
      await screen.findByRole('heading', { name: /your brand, animated\./i }),
    ).toBeInTheDocument()
    expect(screen.getByTestId('location-display')).toHaveTextContent('/')
  })

  it('navigates to the avatar editor when an avatar card is clicked', async () => {
    const user = userEvent.setup()

    server.use(
      http.get(`${API_BASE_URL}/auth/users/me`, () => {
        return HttpResponse.json({
          user: {
            avatarUrl: 'https://avatar.example.com/nico.png',
            email: 'nico@example.com',
            id: 'user-v7',
            name: 'Nico',
          },
        })
      }),
      http.get(`${API_BASE_URL}/creative-studio/avatars`, () => {
        return HttpResponse.json({
          avatars: [{ id: 'avatar-1', name: 'Studio Hero' }],
        })
      }),
      http.get(`${API_BASE_URL}/creative-studio/avatar_configs/:avatarId`, () =>
        HttpResponse.json({ avatar_config: null }),
      ),
    )

    renderApp('/creative-studio')

    await user.click(
      await screen.findByRole('button', {
        name: /open studio hero avatar editor/i,
      }),
    )

    expect(screen.getByTestId('location-display')).toHaveTextContent(
      '/creative-studio/avatars/avatar-1/avatar',
    )
    expect(
      await screen.findByRole('heading', {
        name: /shape the avatar foundation/i,
      }),
    ).toBeInTheDocument()
  })

  it('renders default draft values and image options when avatar config is still empty', async () => {
    server.use(
      http.get(`${API_BASE_URL}/auth/users/me`, () => {
        return HttpResponse.json({
          user: {
            avatarUrl: 'https://avatar.example.com/nico.png',
            email: 'nico@example.com',
            id: 'user-v7',
            name: 'Nico',
          },
        })
      }),
      http.get(`${API_BASE_URL}/creative-studio/avatar_configs/:avatarId`, () =>
        HttpResponse.json({ avatar_config: null }),
      ),
    )

    renderApp('/creative-studio/avatars/avatar-1/avatar')

    expect(await screen.findByLabelText(/avatar description/i)).toHaveValue('')
    expect(screen.getByRole('button', { name: /^2d$/i })).toHaveAttribute(
      'aria-pressed',
      'true',
    )
    expect(
      screen.getByText(/choose a generated direction/i),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: /regenerate/i }),
    ).toBeInTheDocument()
    expect(screen.getAllByRole('radio')).toHaveLength(4)
    expect(
      screen.getByRole('button', { name: /save avatar draft/i }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: /reset draft/i }),
    ).toBeInTheDocument()
  })

  it('lets the user select a generated image option in local ui state', async () => {
    const user = userEvent.setup()

    server.use(
      http.get(`${API_BASE_URL}/auth/users/me`, () => {
        return HttpResponse.json({
          user: {
            avatarUrl: 'https://avatar.example.com/nico.png',
            email: 'nico@example.com',
            id: 'user-v7',
            name: 'Nico',
          },
        })
      }),
      http.get(`${API_BASE_URL}/creative-studio/avatar_configs/:avatarId`, () =>
        HttpResponse.json({ avatar_config: null }),
      ),
    )

    renderApp('/creative-studio/avatars/avatar-1/avatar')

    const radios = (await screen.findAllByRole('radio')) as HTMLInputElement[]

    expect(radios[0]).toBeChecked()
    expect(radios[1]).not.toBeChecked()

    await user.click(radios[1])

    expect(radios[0]).not.toBeChecked()
    expect(radios[1]).toBeChecked()
  })

  it('regenerates image options locally and resets the selected option', async () => {
    const user = userEvent.setup()

    server.use(
      http.get(`${API_BASE_URL}/auth/users/me`, () => {
        return HttpResponse.json({
          user: {
            avatarUrl: 'https://avatar.example.com/nico.png',
            email: 'nico@example.com',
            id: 'user-v7',
            name: 'Nico',
          },
        })
      }),
      http.get(`${API_BASE_URL}/creative-studio/avatar_configs/:avatarId`, () =>
        HttpResponse.json({ avatar_config: null }),
      ),
    )

    renderApp('/creative-studio/avatars/avatar-1/avatar')

    const initialRadios = (await screen.findAllByRole(
      'radio',
    )) as HTMLInputElement[]
    await user.click(initialRadios[2])
    expect(initialRadios[2]).toBeChecked()
    expect(initialRadios[0]).toHaveAttribute('value', '2d-0-0')

    await user.click(screen.getByRole('button', { name: /regenerate/i }))

    const regeneratedRadios = screen.getAllByRole('radio') as HTMLInputElement[]
    expect(regeneratedRadios[0]).toHaveAttribute('value', '2d-1-0')
    expect(regeneratedRadios[0]).toBeChecked()
    expect(regeneratedRadios[2]).not.toBeChecked()
  })

  it('hydrates the avatar editor when a saved draft already exists', async () => {
    server.use(
      http.get(`${API_BASE_URL}/auth/users/me`, () => {
        return HttpResponse.json({
          user: {
            avatarUrl: 'https://avatar.example.com/nico.png',
            email: 'nico@example.com',
            id: 'user-v7',
            name: 'Nico',
          },
        })
      }),
      http.get(`${API_BASE_URL}/creative-studio/avatar_configs/:avatarId`, () =>
        HttpResponse.json({
          avatar_config: {
            avatarId: 'avatar-1',
            artisticStyle: '3D',
            prompt: 'Bold editorial mascot',
          },
        }),
      ),
    )

    renderApp('/creative-studio/avatars/avatar-1/avatar')

    expect(await screen.findByLabelText(/avatar description/i)).toHaveValue(
      'Bold editorial mascot',
    )
    expect(screen.getByRole('button', { name: /^3d$/i })).toHaveAttribute(
      'aria-pressed',
      'true',
    )
  })

  it('saves avatar draft updates successfully', async () => {
    const user = userEvent.setup()
    let savedBody: null | { artisticStyle: string; prompt: string } = null

    server.use(
      http.get(`${API_BASE_URL}/auth/users/me`, () => {
        return HttpResponse.json({
          user: {
            avatarUrl: 'https://avatar.example.com/nico.png',
            email: 'nico@example.com',
            id: 'user-v7',
            name: 'Nico',
          },
        })
      }),
      http.get(`${API_BASE_URL}/creative-studio/avatar_configs/:avatarId`, () =>
        HttpResponse.json({ avatar_config: null }),
      ),
      http.put(
        `${API_BASE_URL}/creative-studio/avatar_configs/:avatarId`,
        async ({ request }) => {
          savedBody = (await request.json()) as {
            artisticStyle: string
            prompt: string
          }

          return HttpResponse.json({
            avatar_config: {
              avatarId: 'avatar-1',
              ...savedBody,
            },
          })
        },
      ),
    )

    renderApp('/creative-studio/avatars/avatar-1/avatar')

    await user.type(
      await screen.findByLabelText(/avatar description/i),
      'Energetic coral storyteller',
    )
    await user.click(screen.getByRole('button', { name: /^3d$/i }))
    await user.click(screen.getByRole('button', { name: /save avatar draft/i }))

    expect(await screen.findByText(/avatar draft saved/i)).toBeInTheDocument()
    expect(savedBody).toEqual({
      artisticStyle: '3D',
      prompt: 'Energetic coral storyteller',
    })
  })

  it('shows avatar draft save errors without leaving the editor', async () => {
    const user = userEvent.setup()

    server.use(
      http.get(`${API_BASE_URL}/auth/users/me`, () => {
        return HttpResponse.json({
          user: {
            avatarUrl: 'https://avatar.example.com/nico.png',
            email: 'nico@example.com',
            id: 'user-v7',
            name: 'Nico',
          },
        })
      }),
      http.get(`${API_BASE_URL}/creative-studio/avatar_configs/:avatarId`, () =>
        HttpResponse.json({ avatar_config: null }),
      ),
      http.put(
        `${API_BASE_URL}/creative-studio/avatar_configs/:avatarId`,
        () => {
          return HttpResponse.json({ message: 'boom' }, { status: 500 })
        },
      ),
    )

    renderApp('/creative-studio/avatars/avatar-1/avatar')

    await user.type(
      await screen.findByLabelText(/avatar description/i),
      'Energetic coral storyteller',
    )
    await user.click(screen.getByRole('button', { name: /save avatar draft/i }))

    expect(
      await screen.findByText(/we could not save your avatar draft/i),
    ).toBeInTheDocument()
    expect(screen.getByTestId('location-display')).toHaveTextContent(
      '/creative-studio/avatars/avatar-1/avatar',
    )
  })

  it('renders placeholder editor steps for future workflow stages', async () => {
    server.use(
      http.get(`${API_BASE_URL}/auth/users/me`, () => {
        return HttpResponse.json({
          user: {
            avatarUrl: 'https://avatar.example.com/nico.png',
            email: 'nico@example.com',
            id: 'user-v7',
            name: 'Nico',
          },
        })
      }),
    )

    renderApp('/creative-studio/avatars/avatar-1/voice')

    expect(
      await screen.findByRole('heading', { name: /voice step/i }),
    ).toBeInTheDocument()
    expect(screen.getByText(/voice will render here soon/i)).toBeInTheDocument()
  })
})
