package http

import (
	"context"
	"fmt"
	stdhttp "net/http"
	"net/url"
	"time"

	authdomain "brandtoonapi/bounded_contexts/identity/auth/domain"
	authusecases "brandtoonapi/bounded_contexts/identity/auth/useCases"
	sessiondomain "brandtoonapi/bounded_contexts/identity/session/domain"
	sessionusecases "brandtoonapi/bounded_contexts/identity/session/useCases"
	userdomain "brandtoonapi/bounded_contexts/identity/user/domain"
	userusecases "brandtoonapi/bounded_contexts/identity/user/useCases"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	sharedconfig "brandtoonapi/bounded_contexts/shared/infra/config"

	"github.com/danielgtaylor/huma/v2"
)

type RouteDependencies struct {
	Config         sharedconfig.Config
	GoogleProvider authdomain.IdentityProvider
	IDGenerator    shareddomain.IDGenerator
	Now            func() time.Time
	SessionRepo    sessiondomain.SessionRepository
	StateCodec     authdomain.OAuthStateCodec
	UserRepo       userdomain.UserRepository
}

type currentUserOutput struct {
	Body struct {
		User currentUserPayload `json:"user"`
	}
}

type currentUserPayload struct {
	AvatarURL string `json:"avatarUrl"`
	Email     string `json:"email"`
	ID        string `json:"id"`
	Name      string `json:"name"`
}

type logoutOutput struct {
	SetCookie stdhttp.Cookie `header:"Set-Cookie"`
	Body      struct {
		Message string `json:"message"`
	}
}

func buildGoogleLoginHandler(deps RouteDependencies) stdhttp.HandlerFunc {
	return func(writer stdhttp.ResponseWriter, request *stdhttp.Request) {
		result, err := authusecases.GetAuthURL(
			authusecases.GetAuthURLQuery{RedirectTo: request.URL.Query().Get("redirectTo")},
			deps.StateCodec,
			deps.GoogleProvider,
			deps.Now(),
		)
		if err != nil {
			redirectToLoginError(writer, request, deps.Config, "oauth_state_error")
			return
		}

		stdhttp.Redirect(writer, request, result.URL, stdhttp.StatusTemporaryRedirect)
	}
}

func buildGoogleCallbackHandler(deps RouteDependencies) stdhttp.HandlerFunc {
	return func(writer stdhttp.ResponseWriter, request *stdhttp.Request) {
		result, err := authusecases.AuthenticateCallback(
			request.Context(),
			authusecases.AuthenticateCallbackCommand{
				Code:       request.URL.Query().Get("code"),
				SessionTTL: deps.Config.SessionTTL,
				State:      request.URL.Query().Get("state"),
			},
			deps.GoogleProvider,
			deps.StateCodec,
			deps.UserRepo,
			deps.SessionRepo,
			deps.IDGenerator,
			deps.Now,
		)
		if err != nil {
			redirectToLoginError(writer, request, deps.Config, "google_auth_failed")
			return
		}

		cookie := newSessionCookie(deps.Config, result.Session.ID, result.Session.ExpiresAt)
		stdhttp.SetCookie(writer, &cookie)
		redirectURL := fmt.Sprintf("%s%s", deps.Config.FrontendBaseURL, result.RedirectTo)
		stdhttp.Redirect(writer, request, redirectURL, stdhttp.StatusTemporaryRedirect)
	}
}

func buildGetCurrentUserHandler(
	deps RouteDependencies,
) func(ctx context.Context, input *struct{}) (*currentUserOutput, error) {
	return func(ctx context.Context, input *struct{}) (*currentUserOutput, error) {
		userMetadata := ctx.Value(shareddomain.UserMetadataContextKey).(*shareddomain.AuthUserMetadata)
		if userMetadata == nil {
			return nil, huma.Error401Unauthorized("missing or invalid session")
		}

		user, err := userusecases.FindUser(
			ctx,
			userusecases.FindUserQuery{UserId: userMetadata.UserId},
			deps.UserRepo,
		)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, huma.Error401Unauthorized("missing or invalid session")
		}

		response := &currentUserOutput{}
		response.Body.User = currentUserPayload{
			AvatarURL: user.AvatarURL,
			Email:     user.Email,
			ID:        user.ID,
			Name:      user.Name,
		}
		return response, nil
	}
}

func buildLogoutHandler(deps RouteDependencies) func(ctx context.Context, input *struct {
	Session stdhttp.Cookie `cookie:"brandtoon_session_id"`
}) (*logoutOutput, error) {
	return func(ctx context.Context, input *struct {
		Session stdhttp.Cookie `cookie:"brandtoon_session_id"`
	}) (*logoutOutput, error) {
		if err := sessionusecases.DeleteSession(
			ctx,
			sessionusecases.LogoutSessionCommand{SessionID: input.Session.Value},
			deps.SessionRepo,
		); err != nil {
			return nil, err
		}

		response := &logoutOutput{
			SetCookie: expiredSessionCookie(deps.Config),
		}
		response.Body.Message = "Logged out"
		return response, nil
	}
}

func redirectToLoginError(
	writer stdhttp.ResponseWriter,
	request *stdhttp.Request,
	config sharedconfig.Config,
	code string,
) {
	loginURL := fmt.Sprintf("%s/login?error=%s", config.FrontendBaseURL, url.QueryEscape(code))
	stdhttp.Redirect(writer, request, loginURL, stdhttp.StatusTemporaryRedirect)
}
