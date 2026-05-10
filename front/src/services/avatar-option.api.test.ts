// @vitest-environment node

import { afterEach, describe, expect, it, vi } from 'vitest'
import {
  deleteAvatarOptions,
  fetchAvatarOptions,
  generateAvatarOptions,
  selectAvatarOption,
} from './avatar-option.api'

describe('avatar option api', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('requests the dedicated avatar-options read endpoint', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        avatar: { avatarId: 'avatar-v7', avatarOptions: [] },
      }),
    })
    vi.stubGlobal('fetch', fetchMock)

    await fetchAvatarOptions('avatar-v7')

    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining('/creative-studio/avatars/avatar-v7/options'),
      expect.objectContaining({ credentials: 'include' }),
    )
  })

  it('posts generation to the dedicated avatar-options command endpoint', async () => {
    const fetchMock = vi.fn().mockResolvedValue({ ok: true })
    vi.stubGlobal('fetch', fetchMock)

    await generateAvatarOptions('avatar-v7')

    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining(
        '/creative-studio/avatars/avatar-v7/options/generate',
      ),
      expect.objectContaining({ credentials: 'include', method: 'POST' }),
    )
  })

  it('posts selection to the dedicated avatar-options endpoint', async () => {
    const fetchMock = vi.fn().mockResolvedValue({ ok: true })
    vi.stubGlobal('fetch', fetchMock)

    await selectAvatarOption('avatar-v7', 'option-v7')

    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining(
        '/creative-studio/avatars/avatar-v7/options/select',
      ),
      expect.objectContaining({
        body: JSON.stringify({ id: 'option-v7' }),
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        method: 'POST',
      }),
    )
  })

  it('deletes options through the dedicated avatar-options endpoint', async () => {
    const fetchMock = vi.fn().mockResolvedValue({ ok: true })
    vi.stubGlobal('fetch', fetchMock)

    await deleteAvatarOptions('avatar-v7', ['option-v7', 'option-v8'])

    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining('/creative-studio/avatars/avatar-v7/options'),
      expect.objectContaining({
        body: JSON.stringify({ ids: ['option-v7', 'option-v8'] }),
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        method: 'DELETE',
      }),
    )
  })
})
