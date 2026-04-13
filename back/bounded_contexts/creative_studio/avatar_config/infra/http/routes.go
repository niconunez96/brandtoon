package avatarconfighttp

import (
	identityauthhttp "brandtoonapi/bounded_contexts/identity/auth/infra/http"
	stdhttp "net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(api huma.API, router chi.Router, deps RouteDependencies) {
	creativeStudioGroup := huma.NewGroup(api, "/creative-studio")
	creativeStudioGroup.UseMiddleware(identityauthhttp.HumaAuthMiddleware(identityauthhttp.AuthMiddlewareDeps{
		SessionRepo: deps.SessionRepo,
		UserRepo:    deps.UserRepo,
		HumaApi:     api,
	}))

	huma.Get(
		creativeStudioGroup,
		"/avatar_configs/{avatarId}",
		buildGetAvatarConfigHandler(deps),
	)
	huma.Register(creativeStudioGroup, huma.Operation{
		OperationID:   "update-avatar-config",
		Method:        stdhttp.MethodPut,
		Path:          "/avatar_configs/{avatarId}",
		Summary:       "Create or update an avatar config draft",
		DefaultStatus: stdhttp.StatusOK,
	}, buildUpdateAvatarConfigHandler(deps))

	_ = router
}
