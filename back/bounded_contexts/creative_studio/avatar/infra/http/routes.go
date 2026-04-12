package avatarhttp

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

	huma.Get(creativeStudioGroup, "/avatars", buildListAvatarsHandler(deps))
	huma.Register(creativeStudioGroup, huma.Operation{
		OperationID:   "create-avatar",
		Method:        stdhttp.MethodPost,
		Path:          "/avatars",
		Summary:       "Create an avatar",
		DefaultStatus: stdhttp.StatusCreated,
	}, buildCreateAvatarHandler(deps))

	_ = router
}
