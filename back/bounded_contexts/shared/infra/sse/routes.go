package sharedsse

import (
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	sharedhttp "brandtoonapi/bounded_contexts/shared/infra/http"
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
)

type RouteDependencies struct {
	Connector *Connector
}

func RegisterRoutes(
	router chi.Router,
	deps RouteDependencies,
	authMiddleware sharedhttp.Middleware,
) {
	router.With(authMiddleware).Get("/events", func(writer stdhttp.ResponseWriter, request *stdhttp.Request) {
		userMetadata, ok := request.Context().Value(shareddomain.UserMetadataContextKey).(*shareddomain.AuthUserMetadata)
		if !ok || userMetadata == nil || userMetadata.UserId == "" {
			writer.WriteHeader(stdhttp.StatusUnauthorized)
			return
		}

		deps.Connector.Stream(writer, request, userMetadata.UserId)
	})
}
