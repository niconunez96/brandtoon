package avatarusecases

import avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"

type AvatarDetailsDTO struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	AvatarOptions []AvatarOptionDTO `json:"avatarOptions"`
}

func serializeAvatarDetails(avatar avatardomain.Avatar) AvatarDetailsDTO {
	return AvatarDetailsDTO{
		ID:            avatar.ID,
		Name:          avatar.Name,
		AvatarOptions: serializeAvatarOptions(avatar).AvatarOptions,
	}
}
