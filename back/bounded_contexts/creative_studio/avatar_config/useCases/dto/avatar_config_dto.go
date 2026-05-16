package avatarconfigdto

import (
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
)

type AvatarConfigDTO struct {
	AvatarID      string `json:"avatarId"`
	ArtisticStyle string `json:"artisticStyle"`
	Personality   string `json:"personality"`
	Prompt        string `json:"prompt"`
}

func Serialize(
	avatarConfig avatarconfigdomain.AvatarConfig,
) AvatarConfigDTO {
	return AvatarConfigDTO{
		AvatarID:      avatarConfig.AvatarID,
		ArtisticStyle: string(avatarConfig.ArtisticStyle),
		Personality:   string(avatarConfig.Personality),
		Prompt:        avatarConfig.Prompt,
	}
}
