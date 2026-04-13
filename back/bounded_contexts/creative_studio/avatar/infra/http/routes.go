package avatarhttp

import (
	sharedhttp "brandtoonapi/bounded_contexts/shared/infra/http"
	stdhttp "net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(
	api huma.API,
	deps RouteDependencies,
	humaMiddlewares ...sharedhttp.HumaMiddleware,
) {
	creativeStudioGroup := huma.NewGroup(api, "/creative-studio")
	creativeStudioGroup.UseMiddleware(humaMiddlewares...)

	huma.Get(creativeStudioGroup, "/avatars", buildListAvatarsHandler(deps))
	huma.Register(creativeStudioGroup, huma.Operation{
		OperationID:   "create-avatar",
		Method:        stdhttp.MethodPost,
		Path:          "/avatars",
		Summary:       "Create an avatar",
		DefaultStatus: stdhttp.StatusCreated,
	}, buildCreateAvatarHandler(deps))
}
