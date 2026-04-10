package http

import (
	stdhttp "net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(api huma.API, router chi.Router, deps RouteDependencies) {
	router.Route("/auth/google", func(r chi.Router) {
		r.Get("/login", buildGoogleLoginHandler(deps))
		r.Get("/callback", buildGoogleCallbackHandler(deps))
	})

	huma.Get(api, "/auth/me", buildGetCurrentUserHandler(deps))
	huma.Register(api, huma.Operation{
		OperationID: "logout-session",
		Method:      stdhttp.MethodPost,
		Path:        "/auth/logout",
		Summary:     "Log out the active session",
	}, buildLogoutHandler(deps))
}
