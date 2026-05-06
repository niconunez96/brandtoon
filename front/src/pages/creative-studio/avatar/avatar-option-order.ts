import type { AvatarOption } from '../../../services/avatar-config.api'

export function orderAvatarOptionsBySelection(options: AvatarOption[]) {
  const selectedOption = options.find((option) => option.selected)
  if (!selectedOption) {
    return options
  }

  return [
    selectedOption,
    ...options.filter((option) => option.id !== selectedOption.id),
  ]
}
