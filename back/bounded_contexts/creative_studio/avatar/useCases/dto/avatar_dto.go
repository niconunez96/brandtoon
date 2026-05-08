package avatardto

import avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"

type AvatarDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func SerializeAvatar(avatar avatardomain.Avatar) AvatarDTO {
	return AvatarDTO{
		ID:   avatar.ID,
		Name: avatar.Name,
	}
}

func SerializeAvatarList(avatars []avatardomain.Avatar) []AvatarDTO {
	avatarDTOs := make([]AvatarDTO, 0, len(avatars))
	for _, avatar := range avatars {
		avatarDTOs = append(avatarDTOs, SerializeAvatar(avatar))
	}

	return avatarDTOs
}
