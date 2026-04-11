package sharedhttp

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type HumaMiddleware func(ctx huma.Context, next func(huma.Context))
type Middleware func(next http.Handler) http.Handler
