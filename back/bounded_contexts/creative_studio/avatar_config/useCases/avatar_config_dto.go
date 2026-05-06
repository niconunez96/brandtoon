package avatarconfigusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarusecases "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
)

type AvatarConfigDTO struct {
	AvatarID      string            `json:"avatarId"`
	ArtisticStyle string            `json:"artisticStyle"`
	Personality   string            `json:"personality"`
	Prompt        string            `json:"prompt"`
	AvatarOptions []avatarusecases.AvatarOptionDTO `json:"avatarOptions"`
}

func serialize(
	avatarConfig avatarconfigdomain.AvatarConfig,
	avatar avatardomain.Avatar,
) AvatarConfigDTO {
	avatarOptions := make([]avatarusecases.AvatarOptionDTO, 0, len(avatar.AvatarOptions))
	for _, option := range avatar.AvatarOptions {
		avatarOptions = append(avatarOptions, avatarusecases.AvatarOptionDTO{
			ID:       option.ID,
			Href:     option.Href,
			Selected: option.Selected,
		})
	}

	return AvatarConfigDTO{
		AvatarID:      avatarConfig.AvatarID,
		ArtisticStyle: string(avatarConfig.ArtisticStyle),
		Personality:   string(avatarConfig.Personality),
		Prompt:        avatarConfig.Prompt,
		AvatarOptions: avatarOptions,
	}
}
