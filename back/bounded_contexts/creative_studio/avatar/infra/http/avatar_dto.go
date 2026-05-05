package avatarhttp

import (
	avatarusecases "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases"
)

type AvatarDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func serializeAvatarUseCaseDTO(avatar avatarusecases.AvatarDTO) AvatarDTO {
	return AvatarDTO{ID: avatar.ID, Name: avatar.Name}
}
