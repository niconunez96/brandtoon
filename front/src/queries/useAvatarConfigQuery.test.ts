// @vitest-environment node

import { afterEach, describe, expect, it, vi } from 'vitest'

const invalidateQueriesMock = vi.fn()
const useMutationMock = vi.fn(
  (options: Record<string, unknown>) => options,
)

vi.mock('@tanstack/react-query', () => ({
  useMutation: (options: Record<string, unknown>) => useMutationMock(options),
  useQuery: vi.fn(),
  useQueryClient: () => ({
    invalidateQueries: invalidateQueriesMock,
  }),
}))

vi.mock('../services/avatar-config.api', () => ({
  deleteAvatarOptions: vi.fn(),
  fetchAvatarConfig: vi.fn(),
  generateAvatarOptions: vi.fn(),
  selectAvatarOption: vi.fn(),
  updateAvatarConfig: vi.fn(),
}))

describe('useAvatarConfigQuery mutations', () => {
	afterEach(() => {
		invalidateQueriesMock.mockReset()
		useMutationMock.mockClear()
		vi.resetModules()
	})

	it('invalidates the avatar query after selecting an option', async () => {
		const { avatarQueryKey } = await import('./useAvatarQuery')
		const { useSelectAvatarOptionMutation } = await import('./useAvatarConfigQuery')

		const mutation = useSelectAvatarOptionMutation('avatar-v7') as {
			onSuccess: () => Promise<void>
		}

		await mutation.onSuccess()

		expect(invalidateQueriesMock).toHaveBeenCalledWith({
			queryKey: avatarQueryKey('avatar-v7'),
		})
	})

	it('invalidates the avatar query after deleting options', async () => {
		const { avatarQueryKey } = await import('./useAvatarQuery')
		const { useDeleteAvatarOptionsMutation } = await import('./useAvatarConfigQuery')

		const mutation = useDeleteAvatarOptionsMutation('avatar-v7') as {
			onSuccess: () => Promise<void>
		}

		await mutation.onSuccess()

		expect(invalidateQueriesMock).toHaveBeenCalledWith({
			queryKey: avatarQueryKey('avatar-v7'),
		})
	})
})
