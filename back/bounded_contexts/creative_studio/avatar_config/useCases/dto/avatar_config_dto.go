package avatarconfigdto

type AvatarConfigDTO struct {
	AvatarID      string            `json:"avatarId"`
	ArtisticStyle string            `json:"artisticStyle"`
	Personality   string            `json:"personality"`
	Prompt        string            `json:"prompt"`
	AvatarOptions []AvatarOptionDTO `json:"avatarOptions"`
}

type AvatarOptionDTO struct {
	ID       string `json:"id"`
	Href     string `json:"href"`
	Selected bool   `json:"selected"`
}
