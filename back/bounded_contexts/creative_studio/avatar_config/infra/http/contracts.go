package avatarconfighttp

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
)

type RouteDependencies struct {
	AvatarConfigRepo avatarconfigdomain.AvatarConfigRepository
	AvatarRepo       avatardomain.AvatarRepository
}

type avatarConfigPayload struct {
	AvatarID      string `json:"avatarId"`
	ArtisticStyle string `json:"artisticStyle"`
	Prompt        string `json:"prompt"`
}

type avatarConfigOutput struct {
	Body struct {
		AvatarConfig *avatarConfigPayload `json:"avatar_config"`
	}
}

type avatarConfigPathInput struct {
	AvatarID string `path:"avatarId"`
}

type updateAvatarConfigBody struct {
	ArtisticStyle string `json:"artisticStyle" enum:"2D,3D"`
	Prompt        string `json:"prompt"                     maxLength:"256"`
}

type updateAvatarConfigInput struct {
	AvatarID string `path:"avatarId"`
	Body     updateAvatarConfigBody
}

func toAvatarConfigPayload(
	avatarConfig avatarconfigdomain.AvatarConfig,
) *avatarConfigPayload {
	return &avatarConfigPayload{
		AvatarID:      avatarConfig.AvatarID,
		ArtisticStyle: string(avatarConfig.ArtisticStyle),
		Prompt:        avatarConfig.Prompt,
	}
}
