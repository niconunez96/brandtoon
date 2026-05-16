package avatardto

import avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"

type AvatarDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func Serialize(avatar avatardomain.Avatar) AvatarDTO {
	return AvatarDTO{
		ID:   avatar.ID,
		Name: avatar.Name,
	}
}

func SerializeList(avatars []avatardomain.Avatar) []AvatarDTO {
	avatarDTOs := make([]AvatarDTO, 0, len(avatars))
	for _, avatar := range avatars {
		avatarDTOs = append(avatarDTOs, Serialize(avatar))
	}

	return avatarDTOs
}
