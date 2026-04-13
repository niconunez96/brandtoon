package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authdomain "brandtoonapi/bounded_contexts/identity/auth/domain"
	authmocks "brandtoonapi/bounded_contexts/identity/auth/domain/mocks"
	authhttp "brandtoonapi/bounded_contexts/identity/auth/infra/http"
	sessionmocks "brandtoonapi/bounded_contexts/identity/session/domain/mocks"
	userdomain "brandtoonapi/bounded_contexts/identity/user/domain"
	usermocks "brandtoonapi/bounded_contexts/identity/user/domain/mocks"
	sharedconfig "brandtoonapi/bounded_contexts/shared/infra/config"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func TestAuthMeReturnsUnauthorizedWithoutCookie(t *testing.T) {
	t.Parallel()

	server := newTestServer(authhttp.RouteDependencies{
		Config: sharedconfig.Config{
			FrontendBaseURL: "http://localhost:5173",
			SessionTTL:      30 * 24 * time.Hour,
		},
		Now:         time.Now,
		SessionRepo: &sessionmocks.SessionRepositoryMock{},
		UserRepo:    &usermocks.UserRepositoryMock{},
	})

	request := httptest.NewRequest(http.MethodGet, "/auth/users/me", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestAuthLogoutClearsCookie(t *testing.T) {
	t.Parallel()

	deletedSessionID := ""
	server := newTestServer(authhttp.RouteDependencies{
		Config: sharedconfig.Config{
			FrontendBaseURL: "http://localhost:5173",
			SessionTTL:      30 * 24 * time.Hour,
		},
		Now: time.Now,
		SessionRepo: &sessionmocks.SessionRepositoryMock{
			DeleteFunc: func(ctx context.Context, id string) error {
				deletedSessionID = id
				return nil
			},
		},
		UserRepo: &usermocks.UserRepositoryMock{},
	})

	request := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	request.AddCookie(&http.Cookie{Name: "brandtoon_session_id", Value: "session-v7"})
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	if deletedSessionID != "session-v7" {
		t.Fatalf("expected session-v7, got %s", deletedSessionID)
	}

	if cookie := recorder.Header().Get("Set-Cookie"); cookie == "" {
		t.Fatal("expected cleared session cookie")
	}
}

func TestGoogleCallbackSetsSessionCookieAndRedirects(t *testing.T) {
	t.Parallel()

	server := newTestServer(authhttp.RouteDependencies{
		Config: sharedconfig.Config{
			AuthStateSecret:    "secret",
			FrontendBaseURL:    "http://localhost:5173",
			GoogleClientID:     "client-id",
			GoogleClientSecret: "client-secret",
			GoogleRedirectURL:  "http://localhost:8888/auth/google/callback",
			SessionTTL:         30 * 24 * time.Hour,
		},
		GoogleProvider: &authmocks.IdentityProviderMock{
			ExchangeCodeFunc: func(ctx context.Context, code string) (*authdomain.Identity, error) {
				return &authdomain.Identity{
					AvatarURL: "https://avatar",
					Email:     "nico@example.com",
					Name:      "Nico",
					Subject:   "google-123",
				}, nil
			},
		},
		IDGenerator: func() (string, error) {
			return "generated-v7", nil
		},
		Now: func() time.Time {
			return time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC)
		},
		SessionRepo: &sessionmocks.SessionRepositoryMock{},
		StateCodec: &authmocks.OAuthStateCodecMock{
			DecodeFunc: func(rawState string) (*authdomain.OAuthState, error) {
				return &authdomain.OAuthState{RedirectTo: "/creative-studio"}, nil
			},
		},
		UserRepo: &usermocks.UserRepositoryMock{
			FindByEmailFunc: func(ctx context.Context, googleSubject string) (*userdomain.User, error) {
				return nil, nil
			},
		},
	})

	request := httptest.NewRequest(
		http.MethodGet,
		"/auth/google/callback?code=auth-code&state=signed-state",
		nil,
	)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected 307, got %d", recorder.Code)
	}

	if location := recorder.Header().
		Get("Location"); location != "http://localhost:5173/creative-studio" {
		t.Fatalf("unexpected redirect location %s", location)
	}

	if cookie := recorder.Header().Get("Set-Cookie"); cookie == "" {
		t.Fatal("expected session cookie to be set")
	}
}

func newTestServer(deps authhttp.RouteDependencies) http.Handler {
	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig("Test API", "1.0.0"))
	authMiddleware := authhttp.HumaAuthMiddleware(authhttp.AuthMiddlewareDeps{
		SessionRepo: deps.SessionRepo,
		UserRepo:    deps.UserRepo,
		HumaApi:     api,
	})
	authhttp.RegisterRoutes(api, router, deps, authMiddleware)
	return router
}
