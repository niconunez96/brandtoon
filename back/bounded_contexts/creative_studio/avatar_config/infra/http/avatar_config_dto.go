package avatarconfighttp

import avatarconfigusecases "brandtoonapi/bounded_contexts/creative_studio/avatar_config/useCases"

type AvatarOptionDTO struct {
	Href     string `json:"href"`
	Selected bool   `json:"selected"`
}

type AvatarConfigDTO struct {
	AvatarID      string            `json:"avatarId"`
	ArtisticStyle string            `json:"artisticStyle"`
	AvatarOptions []AvatarOptionDTO `json:"avatarOptions"`
	Personality   string            `json:"personality"`
	Prompt        string            `json:"prompt"`
}

func serializeAvatarConfigDTO(
	avatarConfig avatarconfigusecases.AvatarConfigDTO,
) *AvatarConfigDTO {
	avatarOptions := make([]AvatarOptionDTO, 0, len(avatarConfig.AvatarOptions))
	for _, option := range avatarConfig.AvatarOptions {
		avatarOptions = append(avatarOptions, AvatarOptionDTO{
			Href:     option.Href,
			Selected: option.Selected,
		})
	}

	return &AvatarConfigDTO{
		AvatarID:      avatarConfig.AvatarID,
		ArtisticStyle: string(avatarConfig.ArtisticStyle),
		AvatarOptions: avatarOptions,
		Personality:   string(avatarConfig.Personality),
		Prompt:        avatarConfig.Prompt,
	}
}
