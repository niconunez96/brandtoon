package avatarhttp

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarusecases "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	identityauthhttp "brandtoonapi/bounded_contexts/identity/auth/infra/http"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
)

type RouteDependencies struct {
	AvatarConfigRepo avatarconfigdomain.AvatarConfigRepository
	AvatarRepo       avatardomain.AvatarRepository
	EventBus         shareddomain.EventBus
	IDGenerator      shareddomain.IDGenerator
	AuthDeps         identityauthhttp.AuthMiddlewareDeps
}

type listAvatarsOutput struct {
	Body struct {
		Avatars []avatarusecases.AvatarDTO `json:"avatars"`
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
		Avatar avatarusecases.AvatarDTO `json:"avatar"`
	}
}

type getAvatarInput struct {
	AvatarID string `path:"avatarId"`
}

type getAvatarOutput struct {
	Body struct {
		Avatar avatarusecases.AvatarDetailsDTO `json:"avatar"`
	}
}

type avatarOptionsPathInput struct {
	AvatarID string `path:"avatarId"`
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
		Avatar avatarusecases.AvatarOptionsDTO `json:"avatar"`
	}
}
