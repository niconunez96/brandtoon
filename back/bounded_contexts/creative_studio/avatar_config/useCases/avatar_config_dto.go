package avatarconfigusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
)

type AvatarConfigDTO struct {
	AvatarID      string `json:"avatarId"`
	ArtisticStyle string `json:"artisticStyle"`
	Personality   string `json:"personality"`
	Prompt        string `json:"prompt"`
	AvatarOptions []AvatarOptionDTO `json:"avatarOptions"`
}

type AvatarOptionDTO struct {
	ID       string `json:"id"`
	Href     string `json:"href"`
	Selected bool   `json:"selected"`
}

func serialize(avatarConfig avatarconfigdomain.AvatarConfig, avatar avatardomain.Avatar) AvatarConfigDTO {
	return AvatarConfigDTO{
		AvatarID:      avatarConfig.AvatarID,
		ArtisticStyle: string(avatarConfig.ArtisticStyle),
		Personality:   string(avatarConfig.Personality),
		Prompt:        avatarConfig.Prompt,
		AvatarOptions: serializeAvatarOptions(avatar.WithNormalizedOptions().AvatarOptions),
	}
}

func serializeAvatarOptions(options []avatardomain.AvatarOption) []AvatarOptionDTO {
	if len(options) == 0 {
		return []AvatarOptionDTO{}
	}

	serialized := make([]AvatarOptionDTO, len(options))
	for index, option := range options {
		serialized[index] = AvatarOptionDTO{
			ID:       option.ID,
			Href:     option.Href,
			Selected: option.Selected,
		}
	}

	return serialized
}
