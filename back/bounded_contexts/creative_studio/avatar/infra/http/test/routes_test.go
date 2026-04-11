package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarmocks "brandtoonapi/bounded_contexts/creative_studio/avatar/domain/mocks"
	avatarhttp "brandtoonapi/bounded_contexts/creative_studio/avatar/infra/http"
	sessiondomain "brandtoonapi/bounded_contexts/identity/session/domain"
	sessionmocks "brandtoonapi/bounded_contexts/identity/session/domain/mocks"
	userdomain "brandtoonapi/bounded_contexts/identity/user/domain"
	usermocks "brandtoonapi/bounded_contexts/identity/user/domain/mocks"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func TestCreativeStudioAvatarsReturnsUnauthorizedWithoutAuth(t *testing.T) {
	t.Parallel()

	server := newTestServer(avatarhttp.RouteDependencies{
		AvatarRepo:  &avatarmocks.AvatarRepositoryMock{},
		SessionRepo: &sessionmocks.SessionRepositoryMock{},
		UserRepo:    &usermocks.UserRepositoryMock{},
	})

	request := httptest.NewRequest(http.MethodGet, "/creative-studio/avatars", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestCreativeStudioAvatarsListsAuthenticatedUserAvatars(t *testing.T) {
	t.Parallel()

	server := newAuthenticatedTestServer(t, avatarhttp.RouteDependencies{
		AvatarRepo: &avatarmocks.AvatarRepositoryMock{
			ListByUserIDFunc: func(ctx context.Context, userID string) ([]avatardomain.Avatar, error) {
				if userID != "user-v7" {
					t.Fatalf("expected user-v7, got %s", userID)
				}

				return []avatardomain.Avatar{
					avatardomain.NewAvatar("avatar-2", "user-v7", "Beta"),
					avatardomain.NewAvatar("avatar-1", "user-v7", "Alpha"),
				}, nil
			},
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/creative-studio/avatars", nil)
	request.AddCookie(&http.Cookie{Name: "brandtoon_session_id", Value: "session-v7"})
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var payload struct {
		Avatars []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"avatars"`
	}
	decodeResponse(t, recorder, &payload)

	if len(payload.Avatars) != 2 {
		t.Fatalf("expected 2 avatars, got %d", len(payload.Avatars))
	}

	if payload.Avatars[0].ID != "avatar-2" {
		t.Fatalf("expected newest avatar first, got %s", payload.Avatars[0].ID)
	}
}

func TestCreativeStudioAvatarsRejectsInvalidName(t *testing.T) {
	t.Parallel()

	server := newAuthenticatedTestServer(t, avatarhttp.RouteDependencies{
		AvatarRepo: &avatarmocks.AvatarRepositoryMock{
			CreateFunc: func(ctx context.Context, avatar avatardomain.Avatar) error {
				t.Fatalf("expected validation to fail before use case execution")
				return nil
			},
		},
	})

	request := httptest.NewRequest(
		http.MethodPost,
		"/creative-studio/avatars",
		bytes.NewBufferString(`{"name":"   "}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "brandtoon_session_id", Value: "session-v7"})
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", recorder.Code)
	}
}

func TestCreativeStudioAvatarsRejectsNameLongerThanOneHundredTwentyCharacters(t *testing.T) {
	t.Parallel()

	server := newAuthenticatedTestServer(t, avatarhttp.RouteDependencies{
		AvatarRepo: &avatarmocks.AvatarRepositoryMock{
			CreateFunc: func(ctx context.Context, avatar avatardomain.Avatar) error {
				t.Fatalf("expected validation to fail before use case execution")
				return nil
			},
		},
	})

	request := httptest.NewRequest(
		http.MethodPost,
		"/creative-studio/avatars",
		bytes.NewBufferString(`{"name":"`+strings.Repeat("a", 121)+`"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "brandtoon_session_id", Value: "session-v7"})
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", recorder.Code)
	}
}

func TestCreativeStudioAvatarsCreatesAvatarForAuthenticatedUser(t *testing.T) {
	t.Parallel()

	server := newAuthenticatedTestServer(t, avatarhttp.RouteDependencies{
		AvatarRepo: &avatarmocks.AvatarRepositoryMock{
			CreateFunc: func(ctx context.Context, avatar avatardomain.Avatar) error {
				if avatar.UserID != "user-v7" {
					t.Fatalf("expected user-v7, got %s", avatar.UserID)
				}

				if avatar.Name != "Studio Hero" {
					t.Fatalf("expected normalized name, got %q", avatar.Name)
				}

				return nil
			},
		},
		IDGenerator: func() (string, error) { return "avatar-v7", nil },
	})

	request := httptest.NewRequest(
		http.MethodPost,
		"/creative-studio/avatars",
		bytes.NewBufferString(`{"name":"  Studio Hero  "}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "brandtoon_session_id", Value: "session-v7"})
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", recorder.Code)
	}

	var payload struct {
		Avatar struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"avatar"`
	}
	decodeResponse(t, recorder, &payload)

	if payload.Avatar.ID != "avatar-v7" {
		t.Fatalf("expected avatar-v7, got %s", payload.Avatar.ID)
	}
}

func newAuthenticatedTestServer(
	t *testing.T,
	deps avatarhttp.RouteDependencies,
) http.Handler {
	t.Helper()

	deps.SessionRepo = &sessionmocks.SessionRepositoryMock{
		FindActiveByIDFunc: func(ctx context.Context, id string) (*sessiondomain.Session, error) {
			return &sessiondomain.Session{
				ID:        id,
				UserID:    "user-v7",
				ExpiresAt: time.Now().Add(24 * time.Hour),
			}, nil
		},
	}
	deps.UserRepo = &usermocks.UserRepositoryMock{
		FindByIDFunc: func(ctx context.Context, id string) (*userdomain.User, error) {
			return &userdomain.User{ID: id, Email: "nico@example.com", Name: "Nico"}, nil
		},
	}

	if deps.IDGenerator == nil {
		deps.IDGenerator = func() (string, error) { return "avatar-generated", nil }
	}

	return newTestServer(deps)
}

func newTestServer(deps avatarhttp.RouteDependencies) http.Handler {
	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig("Test API", "1.0.0"))
	avatarhttp.RegisterRoutes(api, router, deps)
	return router
}

func decodeResponse(t *testing.T, recorder *httptest.ResponseRecorder, target any) {
	t.Helper()

	if err := json.Unmarshal(recorder.Body.Bytes(), target); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}
