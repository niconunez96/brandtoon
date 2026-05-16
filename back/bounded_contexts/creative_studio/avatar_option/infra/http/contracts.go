package avataroptionhttp

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	avataroptiondomain "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain"
	avataroptiondto "brandtoonapi/bounded_contexts/creative_studio/avatar_option/useCases/dto"
	identityauthhttp "brandtoonapi/bounded_contexts/identity/auth/infra/http"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
)

type RouteDependencies struct {
	AvatarConfigRepo avatarconfigdomain.AvatarConfigRepository
	AvatarGenerator  avataroptiondomain.AvatarGenerator
	AvatarRepo       avatardomain.AvatarRepository
	AvatarOptionRepo avataroptiondomain.AvatarOptionRepository
	EventBus         shareddomain.EventBus
	FileStorage      shareddomain.FileStorage
	AuthDeps         identityauthhttp.AuthMiddlewareDeps
}

type listAvatarOptionsInput struct {
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
		Avatar avataroptiondto.AvatarOptionsDTO `json:"avatar"`
	}
}
