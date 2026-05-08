package avatarusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatardto "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases/dto"
)

func serializeAvatarOptions(avatar avatardomain.Avatar) avatardto.AvatarOptionsDTO {
	avatarOptions := make([]avatardto.AvatarOptionDTO, 0, len(avatar.AvatarOptions))
	for _, option := range avatar.AvatarOptions {
		avatarOptions = append(avatarOptions, avatardto.AvatarOptionDTO{
			ID:       option.ID,
			Href:     option.Href,
			Selected: option.Selected,
		})
	}

	return avatardto.AvatarOptionsDTO{
		AvatarID:      avatar.ID,
		AvatarOptions: avatarOptions,
	}
}
