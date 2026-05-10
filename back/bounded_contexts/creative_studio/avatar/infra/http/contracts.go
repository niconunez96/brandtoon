package avatarhttp

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatardto "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases/dto"
	identityauthhttp "brandtoonapi/bounded_contexts/identity/auth/infra/http"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
)

type RouteDependencies struct {
	AvatarRepo  avatardomain.AvatarRepository
	IDGenerator shareddomain.IDGenerator
	AuthDeps    identityauthhttp.AuthMiddlewareDeps
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
		Avatar avatardto.AvatarDTO `json:"avatar"`
	}
}
