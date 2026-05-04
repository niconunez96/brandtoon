package avatarconfigusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
)

type AvatarConfigDTO struct {
	AvatarID      string
	ArtisticStyle string
	Personality   string
	Prompt        string
	AvatarOptions []avatardomain.AvatarOption
}

func serialize(
	avatarConfig avatarconfigdomain.AvatarConfig,
	avatar avatardomain.Avatar,
) AvatarConfigDTO {
	return AvatarConfigDTO{
		AvatarID:      avatarConfig.AvatarID,
		ArtisticStyle: string(avatarConfig.ArtisticStyle),
		Personality:   string(avatarConfig.Personality),
		Prompt:        avatarConfig.Prompt,
		AvatarOptions: avatar.AvatarOptions,
	}
}
