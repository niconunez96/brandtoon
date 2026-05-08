package avatardto

type AvatarOptionDTO struct {
	ID       string `json:"id"`
	Href     string `json:"href"`
	Selected bool   `json:"selected"`
}

type AvatarOptionsDTO struct {
	AvatarID      string            `json:"avatarId"`
	AvatarOptions []AvatarOptionDTO `json:"avatarOptions"`
}
