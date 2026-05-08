package avatarconfighttp

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	avatarconfigdto "brandtoonapi/bounded_contexts/creative_studio/avatar_config/useCases/dto"
)

type RouteDependencies struct {
	AvatarConfigRepo avatarconfigdomain.AvatarConfigRepository
	AvatarRepo       avatardomain.AvatarRepository
}

type avatarConfigOutput struct {
	Body struct {
		AvatarConfig *avatarconfigdto.AvatarConfigDTO `json:"avatar_config"`
	}
}

type avatarConfigPathInput struct {
	AvatarID string `path:"avatarId"`
}

type updateAvatarConfigBody struct {
	ArtisticStyle string `json:"artisticStyle" enum:"2D,3D"`
	Personality   string `json:"personality"   enum:"Friendly,Bold,Playful"`
	Prompt        string `json:"prompt"                                     maxLength:"256"`
}

type updateAvatarConfigInput struct {
	AvatarID string `path:"avatarId"`
	Body     updateAvatarConfigBody
}
