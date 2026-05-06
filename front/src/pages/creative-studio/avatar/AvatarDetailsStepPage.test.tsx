import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { afterEach, describe, expect, it } from 'vitest'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { API_BASE_URL } from '../../../shared/config/api'
import { server } from '../../../test/mocks/server'
import { AvatarDetailsStepPage } from './AvatarDetailsStepPage'

afterEach(() => {
  server.resetHandlers()
})

function renderAvatarDetailsPage() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  })

  return render(
    <MemoryRouter initialEntries={['/creative-studio/avatars/avatar-1/avatar']}>
      <QueryClientProvider client={queryClient}>
        <Routes>
          <Route
            element={<AvatarDetailsStepPage />}
            path="/creative-studio/avatars/:avatarId/avatar"
          />
        </Routes>
      </QueryClientProvider>
    </MemoryRouter>,
  )
}

describe('AvatarDetailsStepPage', () => {
  it('selects avatar options and deletes multiple selected options', async () => {
    const user = userEvent.setup()
    const selectBodies: Array<{ id: string }> = []
    const deleteBodies: Array<{ ids: string[] }> = []
    let avatarConfig = {
      avatarId: 'avatar-1',
      artisticStyle: '2D',
      avatarOptions: [
        {
          id: 'option-1',
          href: 'https://cdn.brandtoon.local/avatars/avatar-1/options/1.png',
          selected: false,
        },
        {
          id: 'option-2',
          href: 'https://cdn.brandtoon.local/avatars/avatar-1/options/2.png',
          selected: false,
        },
        {
          id: 'option-3',
          href: 'https://cdn.brandtoon.local/avatars/avatar-1/options/3.png',
          selected: false,
        },
      ],
      personality: 'Friendly',
      prompt: 'Studio mascot',
    }

    server.use(
      http.get(`${API_BASE_URL}/creative-studio/avatar_configs/:avatarId`, () =>
        HttpResponse.json({ avatar_config: avatarConfig }),
      ),
      http.post(
        `${API_BASE_URL}/creative-studio/avatar_configs/:avatarId/options/select`,
        async ({ request }) => {
          const body = (await request.json()) as { id: string }
          selectBodies.push(body)
          avatarConfig = {
            ...avatarConfig,
            avatarOptions: avatarConfig.avatarOptions.map((option) => ({
              ...option,
              selected: option.id === body.id,
            })),
          }

          return HttpResponse.json({ avatar: avatarConfig })
        },
      ),
      http.delete(
        `${API_BASE_URL}/creative-studio/avatar_configs/:avatarId/options`,
        async ({ request }) => {
          const body = (await request.json()) as { ids: string[] }
          deleteBodies.push(body)
          avatarConfig = {
            ...avatarConfig,
            avatarOptions: avatarConfig.avatarOptions.filter(
              (option) => !body.ids.includes(option.id),
            ),
          }

          return HttpResponse.json({ avatar: avatarConfig })
        },
      ),
    )

    renderAvatarDetailsPage()

    await user.click(
      await screen.findByRole('button', {
        name: /select avatar option option-2/i,
      }),
    )

    await waitFor(() =>
      expect(selectBodies).toEqual([{ id: 'option-2' }]),
    )
    expect(
      screen.getByRole('button', { name: /delete selected \(0\)/i }),
    ).toBeInTheDocument()
    expect(screen.getByText('Selected', { selector: 'span' })).toBeInTheDocument()

    await user.click(
      screen.getByLabelText(/mark avatar option option-1 for deletion/i),
    )
    await user.click(
      screen.getByLabelText(/mark avatar option option-3 for deletion/i),
    )
    await user.click(
      screen.getByRole('button', { name: /delete selected \(2\)/i }),
    )

    await waitFor(() =>
      expect(deleteBodies).toEqual([{ ids: ['option-1', 'option-3'] }]),
    )
    expect(
      screen.queryByLabelText(/mark avatar option option-1 for deletion/i),
    ).not.toBeInTheDocument()
    expect(
      screen.queryByLabelText(/mark avatar option option-3 for deletion/i),
    ).not.toBeInTheDocument()
    expect(
      screen.getByLabelText(/mark avatar option option-2 for deletion/i),
    ).toBeInTheDocument()
  })
})
