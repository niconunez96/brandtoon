package test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarmocks "brandtoonapi/bounded_contexts/creative_studio/avatar/domain/mocks"
	avataroptiondomain "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain"
	avataroptionmocks "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain/mocks"
	avataroptionhttp "brandtoonapi/bounded_contexts/creative_studio/avatar_option/infra/http"
	authhttp "brandtoonapi/bounded_contexts/identity/auth/infra/http"
	sessiondomain "brandtoonapi/bounded_contexts/identity/session/domain"
	sessionmocks "brandtoonapi/bounded_contexts/identity/session/domain/mocks"
	userdomain "brandtoonapi/bounded_contexts/identity/user/domain"
	usermocks "brandtoonapi/bounded_contexts/identity/user/domain/mocks"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func TestAvatarOptionsListRouteReturnsDedicatedCollection(t *testing.T) {
	t.Parallel()

	href := "https://cdn.brandtoon.local/avatars/avatar-v7/options/1.png"
	server := newAuthenticatedTestServer(t, avataroptionhttp.RouteDependencies{
		AvatarRepo: &avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
		},
		AvatarOptionRepo: &avataroptionmocks.AvatarOptionRepositoryMock{
			ListByAvatarIDFunc: func(ctx context.Context, avatarID string) ([]avataroptiondomain.AvatarOption, error) {
				return []avataroptiondomain.AvatarOption{
					{
						ID:       "option-1",
						AvatarID: avatarID,
						Status:   avataroptiondomain.StatusDone,
						Href:     &href,
						Selected: true,
					},
				}, nil
			},
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/creative-studio/avatars/avatar-v7/options", nil)
	request.AddCookie(&http.Cookie{Name: "brandtoon_session_id", Value: "session-v7"})
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestAvatarOptionsSelectRouteUsesDedicatedEndpoint(t *testing.T) {
	t.Parallel()

	server := newAuthenticatedTestServer(t, avataroptionhttp.RouteDependencies{
		AvatarRepo: &avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
		},
		AvatarOptionRepo: &avataroptionmocks.AvatarOptionRepositoryMock{
			SelectFunc: func(ctx context.Context, avatarID string, optionID string) error { return nil },
			ListByAvatarIDFunc: func(ctx context.Context, avatarID string) ([]avataroptiondomain.AvatarOption, error) {
				href := "https://cdn.brandtoon.local/avatars/avatar-v7/options/2.png"
				return []avataroptiondomain.AvatarOption{
					{
						ID:       "option-2",
						AvatarID: avatarID,
						Status:   avataroptiondomain.StatusDone,
						Href:     &href,
						Selected: true,
					},
				}, nil
			},
		},
	})

	request := httptest.NewRequest(
		http.MethodPost,
		"/creative-studio/avatars/avatar-v7/options/select",
		bytes.NewBufferString(`{"id":"option-2"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "brandtoon_session_id", Value: "session-v7"})
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestAvatarOptionsDeleteRouteUsesDedicatedEndpoint(t *testing.T) {
	t.Parallel()

	server := newAuthenticatedTestServer(t, avataroptionhttp.RouteDependencies{
		AvatarRepo: &avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
		},
		AvatarOptionRepo: &avataroptionmocks.AvatarOptionRepositoryMock{
			DeleteFunc: func(ctx context.Context, avatarID string, ids []string) error { return nil },
			ListByAvatarIDFunc: func(ctx context.Context, avatarID string) ([]avataroptiondomain.AvatarOption, error) {
				return []avataroptiondomain.AvatarOption{}, nil
			},
		},
	})

	request := httptest.NewRequest(
		http.MethodDelete,
		"/creative-studio/avatars/avatar-v7/options",
		bytes.NewBufferString(`{"ids":["option-1"]}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "brandtoon_session_id", Value: "session-v7"})
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
}

func newAuthenticatedTestServer(t *testing.T, deps avataroptionhttp.RouteDependencies) http.Handler {
	t.Helper()

	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig("Test API", "1.0.0"))
	plainAuthMiddleware := authhttp.AuthMiddleware(authhttp.AuthMiddlewareDeps{
		SessionRepo: &sessionmocks.SessionRepositoryMock{
			FindActiveByIDFunc: func(ctx context.Context, id string) (*sessiondomain.Session, error) {
				return &sessiondomain.Session{ID: id, UserID: "user-v7", ExpiresAt: time.Now().Add(24 * time.Hour)}, nil
			},
		},
		UserRepo: &usermocks.UserRepositoryMock{
			FindByIDFunc: func(ctx context.Context, id string) (*userdomain.User, error) {
				return &userdomain.User{ID: id, Email: "nico@example.com", Name: "Nico"}, nil
			},
		},
		HumaApi: api,
	})
	authMiddleware := authhttp.HumaAuthMiddleware(authhttp.AuthMiddlewareDeps{
		SessionRepo: &sessionmocks.SessionRepositoryMock{
			FindActiveByIDFunc: func(ctx context.Context, id string) (*sessiondomain.Session, error) {
				return &sessiondomain.Session{ID: id, UserID: "user-v7", ExpiresAt: time.Now().Add(24 * time.Hour)}, nil
			},
		},
		UserRepo: &usermocks.UserRepositoryMock{
			FindByIDFunc: func(ctx context.Context, id string) (*userdomain.User, error) {
				return &userdomain.User{ID: id, Email: "nico@example.com", Name: "Nico"}, nil
			},
		},
		HumaApi: api,
	})
	avataroptionhttp.RegisterRoutes(api, router, deps, plainAuthMiddleware, authMiddleware)
	return router
}
