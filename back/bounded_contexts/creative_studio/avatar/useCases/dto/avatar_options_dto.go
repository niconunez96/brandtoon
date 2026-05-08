package avatardto

import avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"

type AvatarOptionDTO struct {
	ID       string `json:"id"`
	Href     string `json:"href"`
	Selected bool   `json:"selected"`
}

type AvatarOptionsDTO struct {
	AvatarID      string            `json:"avatarId"`
	AvatarOptions []AvatarOptionDTO `json:"avatarOptions"`
}

func SerializeAvatarOptions(avatar avatardomain.Avatar) AvatarOptionsDTO {
	avatarOptions := make([]AvatarOptionDTO, 0, len(avatar.AvatarOptions))
	for _, option := range avatar.AvatarOptions {
		avatarOptions = append(avatarOptions, AvatarOptionDTO{
			ID:       option.ID,
			Href:     option.Href,
			Selected: option.Selected,
		})
	}

	return AvatarOptionsDTO{
		AvatarID:      avatar.ID,
		AvatarOptions: avatarOptions,
	}
}
