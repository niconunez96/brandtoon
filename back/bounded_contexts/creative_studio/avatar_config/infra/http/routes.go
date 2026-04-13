package avatarconfighttp

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
}
