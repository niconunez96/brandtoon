package avatarusecases

import avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"

type AvatarDTO struct {
	ID   string
	Name string
}

func serialize(avatar avatardomain.Avatar) AvatarDTO {
	return AvatarDTO{
		ID:   avatar.ID,
		Name: avatar.Name,
	}
}

func serializeList(avatars []avatardomain.Avatar) []AvatarDTO {
	avatarDTOs := make([]AvatarDTO, 0, len(avatars))
	for _, avatar := range avatars {
		avatarDTOs = append(avatarDTOs, serialize(avatar))
	}

	return avatarDTOs
}
