package avatardto

type AvatarDetailsDTO struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	AvatarOptions []AvatarOptionDTO `json:"avatarOptions"`
}
