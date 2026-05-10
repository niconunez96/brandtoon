package avataroptiondto

import avataroptiondomain "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain"

type AvatarOptionDTO struct {
	ID       string  `json:"id"`
	Status   string  `json:"status"`
	Selected bool    `json:"selected"`
	Href     *string `json:"href,omitempty"`
}

type AvatarOptionsDTO struct {
	AvatarID      string            `json:"avatarId"`
	AvatarOptions []AvatarOptionDTO `json:"avatarOptions"`
}

func Serialize(option avataroptiondomain.AvatarOption) AvatarOptionDTO {
	return AvatarOptionDTO{
		ID:       option.ID,
		Status:   string(option.Status),
		Selected: option.Selected,
		Href:     option.Href,
	}
}

func SerializeList(avatarID string, options []avataroptiondomain.AvatarOption) AvatarOptionsDTO {
	avatarOptions := make([]AvatarOptionDTO, 0, len(options))
	for _, option := range options {
		avatarOptions = append(avatarOptions, Serialize(option))
	}

	return AvatarOptionsDTO{AvatarID: avatarID, AvatarOptions: avatarOptions}
}
