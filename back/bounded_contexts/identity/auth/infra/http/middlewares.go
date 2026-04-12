package authhttp

import (
	sessiondomain "brandtoonapi/bounded_contexts/identity/session/domain"
	sessionusecases "brandtoonapi/bounded_contexts/identity/session/useCases"
	userdomain "brandtoonapi/bounded_contexts/identity/user/domain"
	userusecases "brandtoonapi/bounded_contexts/identity/user/useCases"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	sharedhttp "brandtoonapi/bounded_contexts/shared/infra/http"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type AuthMiddlewareDeps struct {
	SessionRepo sessiondomain.SessionRepository
	UserRepo    userdomain.UserRepository
	HumaApi     huma.API
}

func AuthMiddleware(deps AuthMiddlewareDeps) sharedhttp.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie(sessionCookieName)
			if err != nil || sessionCookie == nil {
				return
			}
			session, err := sessionusecases.FindSession(
				r.Context(),
				sessionusecases.FindSessionQuery{SessionId: sessionCookie.Value},
				deps.SessionRepo,
			)
			if err != nil || session == nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			user, err := userusecases.FindUser(
				r.Context(),
				userusecases.FindUserQuery{UserId: session.UserID},
				deps.UserRepo,
			)
			if err != nil || user == nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(
				r.Context(),
				shareddomain.UserMetadataContextKey,
				&shareddomain.AuthUserMetadata{UserId: user.ID},
			)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)

		})
	}
}

func HumaAuthMiddleware(deps AuthMiddlewareDeps) sharedhttp.HumaMiddleware {
	return func(humaCtx huma.Context, next func(ctx huma.Context)) {
		sessionCookie, err := huma.ReadCookie(humaCtx, sessionCookieName)
		if err != nil || sessionCookie == nil {
			_ = huma.WriteErr(deps.HumaApi, humaCtx, http.StatusUnauthorized, "unauthorized")
			return
		}
		session, err := sessionusecases.FindSession(
			humaCtx.Context(),
			sessionusecases.FindSessionQuery{SessionId: sessionCookie.Value},
			deps.SessionRepo,
		)
		if err != nil || session == nil {
			_ = huma.WriteErr(deps.HumaApi, humaCtx, http.StatusUnauthorized, "unauthorized")
			return
		}
		user, err := userusecases.FindUser(
			humaCtx.Context(),
			userusecases.FindUserQuery{UserId: session.UserID},
			deps.UserRepo,
		)
		if err != nil || user == nil {
			_ = huma.WriteErr(deps.HumaApi, humaCtx, http.StatusUnauthorized, "unauthorized")
			return
		}
		ctx := context.WithValue(
			humaCtx.Context(),
			shareddomain.UserMetadataContextKey,
			&shareddomain.AuthUserMetadata{UserId: user.ID},
		)
		humaNewCtx := huma.WithContext(humaCtx, ctx)
		next(humaNewCtx)
	}
}
