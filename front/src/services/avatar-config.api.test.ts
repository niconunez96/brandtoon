// @vitest-environment node

import { afterEach, describe, expect, it, vi } from 'vitest'
import { fetchAvatarConfig, updateAvatarConfig } from './avatar-config.api'

describe('avatar config api', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('requests the avatar-config read endpoint', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ avatar_config: null }),
    })
    vi.stubGlobal('fetch', fetchMock)

    await fetchAvatarConfig('avatar-v7')

    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining('/creative-studio/avatar_configs/avatar-v7'),
      expect.objectContaining({
        credentials: 'include',
      }),
    )
  })

  it('updates the avatar-config draft payload', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ avatar_config: null }),
    })
    vi.stubGlobal('fetch', fetchMock)

    await updateAvatarConfig('avatar-v7', {
      artisticStyle: '2D',
      personality: 'Friendly',
      prompt: 'Studio mascot',
    })

    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining('/creative-studio/avatar_configs/avatar-v7'),
      expect.objectContaining({
        body: JSON.stringify({
          artisticStyle: '2D',
          personality: 'Friendly',
          prompt: 'Studio mascot',
        }),
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        method: 'PUT',
      }),
    )
  })
})
