package avataroptionhttp

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

	huma.Get(creativeStudioGroup, "/avatars/{avatarId}/options", buildListAvatarOptionsHandler(deps))
	huma.Register(
		creativeStudioGroup,
		huma.Operation{
			OperationID:   "select-avatar-option",
			Method:        stdhttp.MethodPost,
			Path:          "/avatars/{avatarId}/options/select",
			Summary:       "Select an avatar option",
			DefaultStatus: stdhttp.StatusOK,
		},
		buildSelectAvatarOptionHandler(deps),
	)
	huma.Register(
		creativeStudioGroup,
		huma.Operation{
			OperationID:   "delete-avatar-options",
			Method:        stdhttp.MethodDelete,
			Path:          "/avatars/{avatarId}/options",
			Summary:       "Delete multiple avatar options",
			DefaultStatus: stdhttp.StatusOK,
		},
		buildDeleteAvatarOptionsHandler(deps),
	)
	router.With(authMiddleware).
		Post("/creative-studio/avatars/{avatarId}/options/generate", buildGenerateAvatarOptionsHandler(deps))
}
