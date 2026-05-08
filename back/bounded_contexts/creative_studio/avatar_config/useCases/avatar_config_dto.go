package avatarconfigusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	avatarconfigdto "brandtoonapi/bounded_contexts/creative_studio/avatar_config/useCases/dto"
)

func serialize(
	avatarConfig avatarconfigdomain.AvatarConfig,
	avatar avatardomain.Avatar,
) avatarconfigdto.AvatarConfigDTO {
	return avatarconfigdto.AvatarConfigDTO{
		AvatarID:      avatarConfig.AvatarID,
		ArtisticStyle: string(avatarConfig.ArtisticStyle),
		Personality:   string(avatarConfig.Personality),
		Prompt:        avatarConfig.Prompt,
		AvatarOptions: serializeAvatarOptions(avatar.WithNormalizedOptions().AvatarOptions),
	}
}

func serializeAvatarOptions(options []avatardomain.AvatarOption) []avatarconfigdto.AvatarOptionDTO {
	if len(options) == 0 {
		return []avatarconfigdto.AvatarOptionDTO{}
	}

	serialized := make([]avatarconfigdto.AvatarOptionDTO, len(options))
	for index, option := range options {
		serialized[index] = avatarconfigdto.AvatarOptionDTO{
			ID:       option.ID,
			Href:     option.Href,
			Selected: option.Selected,
		}
	}

	return serialized
}
