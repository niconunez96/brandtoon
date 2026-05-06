import { afterEach, describe, expect, it, vi } from 'vitest'
import { deleteAvatarOptions, selectAvatarOption } from './avatar-config.api'

describe('selectAvatarOption', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('posts the selected option id to the select endpoint', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ avatar_config: null }),
    })
    vi.stubGlobal('fetch', fetchMock)

    await selectAvatarOption('avatar-v7', 'option-v7')

    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining(
        '/creative-studio/avatar_configs/avatar-v7/options/select',
      ),
      expect.objectContaining({
        body: JSON.stringify({ id: 'option-v7' }),
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        method: 'POST',
      }),
    )
  })

  it('posts selected ids to the delete endpoint', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ avatar_config: null }),
    })
    vi.stubGlobal('fetch', fetchMock)

    await deleteAvatarOptions('avatar-v7', ['option-v7', 'option-v8'])

    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining(
        '/creative-studio/avatar_configs/avatar-v7/options',
      ),
      expect.objectContaining({
        body: JSON.stringify({ ids: ['option-v7', 'option-v8'] }),
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        method: 'DELETE',
      }),
    )
  })
})
