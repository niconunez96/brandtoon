package avatarhttp

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatardto "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases/dto"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	identityauthhttp "brandtoonapi/bounded_contexts/identity/auth/infra/http"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
)

type RouteDependencies struct {
	AvatarConfigRepo avatarconfigdomain.AvatarConfigRepository
	AvatarGenerator  avatardomain.AvatarGenerator
	AvatarRepo       avatardomain.AvatarRepository
	EventBus         shareddomain.EventBus
	FileStorage      shareddomain.FileStorage
	IDGenerator      shareddomain.IDGenerator
	AuthDeps         identityauthhttp.AuthMiddlewareDeps
}

type listAvatarsOutput struct {
	Body struct {
		Avatars []avatardto.AvatarDTO `json:"avatars"`
	}
}

type createAvatarBody struct {
	Name string `json:"name" minLength:"1" maxLength:"120" pattern:"\S" patternDescription:"must contain at least one non-whitespace character"` // nolint
}

type createAvatarInput struct {
	Body createAvatarBody
}

type createAvatarOutput struct {
	Body struct {
		Avatar avatardto.AvatarDTO `json:"avatar"`
	}
}

type getAvatarInput struct {
	AvatarID string `path:"avatarId"`
}

type getAvatarOutput struct {
	Body struct {
		Avatar avatardto.AvatarDetailsDTO `json:"avatar"`
	}
}

type selectAvatarOptionBody struct {
	ID string `json:"id" minLength:"1"`
}

type selectAvatarOptionInput struct {
	AvatarID string `path:"avatarId"`
	Body     selectAvatarOptionBody
}

type deleteAvatarOptionsBody struct {
	IDs []string `json:"ids" minItems:"1"`
}

type deleteAvatarOptionsInput struct {
	AvatarID string `path:"avatarId"`
	Body     deleteAvatarOptionsBody
}

type avatarOptionsOutput struct {
	Body struct {
		Avatar avatardto.AvatarOptionsDTO `json:"avatar"`
	}
}
