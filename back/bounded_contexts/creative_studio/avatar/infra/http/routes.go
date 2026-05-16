package avatarhttp

import (
	sharedhttp "brandtoonapi/bounded_contexts/shared/infra/http"
	stdhttp "net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(
	api huma.API,
	router chi.Router,
	deps RouteDependencies,
	authMiddleware sharedhttp.Middleware,
	humaMiddlewares ...sharedhttp.HumaMiddleware,
) {
	creativeStudioGroup := huma.NewGroup(api, "/creative-studio")
	creativeStudioGroup.UseMiddleware(humaMiddlewares...)

	huma.Get(creativeStudioGroup, "/avatars", buildListAvatarsHandler(deps))
	huma.Get(creativeStudioGroup, "/avatars/{avatarId}", buildGetAvatarHandler(deps))
	huma.Register(creativeStudioGroup, huma.Operation{
		OperationID:   "create-avatar",
		Method:        stdhttp.MethodPost,
		Path:          "/avatars",
		Summary:       "Create an avatar",
		DefaultStatus: stdhttp.StatusCreated,
	}, buildCreateAvatarHandler(deps))
}
