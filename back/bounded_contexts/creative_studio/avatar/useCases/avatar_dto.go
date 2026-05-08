package avatarusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatardto "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases/dto"
)

func serialize(avatar avatardomain.Avatar) avatardto.AvatarDTO {
	return avatardto.AvatarDTO{
		ID:   avatar.ID,
		Name: avatar.Name,
	}
}

func serializeList(avatars []avatardomain.Avatar) []avatardto.AvatarDTO {
	avatarDTOs := make([]avatardto.AvatarDTO, 0, len(avatars))
	for _, avatar := range avatars {
		avatarDTOs = append(avatarDTOs, serialize(avatar))
	}

	return avatarDTOs
}
