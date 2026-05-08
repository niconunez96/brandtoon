package avatarusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatardto "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases/dto"
)

func serializeAvatarDetails(avatar avatardomain.Avatar) avatardto.AvatarDetailsDTO {
	return avatardto.AvatarDetailsDTO{
		ID:            avatar.ID,
		Name:          avatar.Name,
		AvatarOptions: serializeAvatarOptions(avatar).AvatarOptions,
	}
}
