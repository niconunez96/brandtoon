// @vitest-environment node

import { afterEach, describe, expect, it, vi } from 'vitest'
import { fetchAvatar } from './avatar.api'

describe('fetchAvatar', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('requests the avatar details endpoint', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ avatar: null }),
    })
    vi.stubGlobal('fetch', fetchMock)

    await fetchAvatar('avatar-v7')

    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining('/creative-studio/avatars/avatar-v7'),
      expect.objectContaining({
        credentials: 'include',
      }),
    )
  })
})
