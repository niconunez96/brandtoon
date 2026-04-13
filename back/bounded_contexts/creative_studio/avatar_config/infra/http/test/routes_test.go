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
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	avatarconfigmocks "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain/mocks"
	avatarconfighttp "brandtoonapi/bounded_contexts/creative_studio/avatar_config/infra/http"
	sessiondomain "brandtoonapi/bounded_contexts/identity/session/domain"
	sessionmocks "brandtoonapi/bounded_contexts/identity/session/domain/mocks"
	userdomain "brandtoonapi/bounded_contexts/identity/user/domain"
	usermocks "brandtoonapi/bounded_contexts/identity/user/domain/mocks"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func TestAvatarConfigReturnsUnauthorizedWithoutAuth(t *testing.T) {
	t.Parallel()

	server := newAvatarConfigTestServer(avatarconfighttp.RouteDependencies{
		AvatarConfigRepo: &avatarconfigmocks.AvatarConfigRepositoryMock{},
		AvatarRepo:       &avatarmocks.AvatarRepositoryMock{},
		SessionRepo:      &sessionmocks.SessionRepositoryMock{},
		UserRepo:         &usermocks.UserRepositoryMock{},
	})

	request := httptest.NewRequest(http.MethodGet, "/creative-studio/avatar_configs/avatar-v7", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestAvatarConfigGetReturnsNullPayloadWhenDraftIsMissing(t *testing.T) {
	t.Parallel()

	server := newAuthenticatedAvatarConfigTestServer(t, avatarconfighttp.RouteDependencies{
		AvatarConfigRepo: &avatarconfigmocks.AvatarConfigRepositoryMock{},
		AvatarRepo: &avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/creative-studio/avatar_configs/avatar-v7", nil)
	request.AddCookie(&http.Cookie{Name: "brandtoon_session_id", Value: "session-v7"})
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var payload struct {
		AvatarConfig any `json:"avatar_config"`
	}
	decodeAvatarConfigResponse(t, recorder, &payload)

	if payload.AvatarConfig != nil {
		t.Fatalf("expected null avatar_config, got %+v", payload.AvatarConfig)
	}
}

func TestAvatarConfigGetReturnsNotFoundForMissingAvatar(t *testing.T) {
	t.Parallel()

	server := newAuthenticatedAvatarConfigTestServer(t, avatarconfighttp.RouteDependencies{
		AvatarConfigRepo: &avatarconfigmocks.AvatarConfigRepositoryMock{},
		AvatarRepo:       &avatarmocks.AvatarRepositoryMock{},
	})

	request := httptest.NewRequest(http.MethodGet, "/creative-studio/avatar_configs/avatar-v7", nil)
	request.AddCookie(&http.Cookie{Name: "brandtoon_session_id", Value: "session-v7"})
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", recorder.Code)
	}
}

func TestAvatarConfigPutRejectsInvalidArtisticStyle(t *testing.T) {
	t.Parallel()

	server := newAuthenticatedAvatarConfigTestServer(t, avatarconfighttp.RouteDependencies{
		AvatarConfigRepo: &avatarconfigmocks.AvatarConfigRepositoryMock{},
		AvatarRepo: &avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
		},
	})

	request := httptest.NewRequest(
		http.MethodPut,
		"/creative-studio/avatar_configs/avatar-v7",
		bytes.NewBufferString(`{"prompt":"hello","artisticStyle":"Clay"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "brandtoon_session_id", Value: "session-v7"})
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", recorder.Code)
	}
}

func TestAvatarConfigPutCreatesOrUpdatesDraft(t *testing.T) {
	t.Parallel()

	persistedConfig := avatarconfigdomain.AvatarConfig{}
	server := newAuthenticatedAvatarConfigTestServer(t, avatarconfighttp.RouteDependencies{
		AvatarConfigRepo: &avatarconfigmocks.AvatarConfigRepositoryMock{
			UpsertFunc: func(ctx context.Context, avatarConfig avatarconfigdomain.AvatarConfig) error {
				persistedConfig = avatarConfig
				return nil
			},
		},
		AvatarRepo: &avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
		},
	})

	request := httptest.NewRequest(
		http.MethodPut,
		"/creative-studio/avatar_configs/avatar-v7",
		bytes.NewBufferString(`{"prompt":"Bold mascot","artisticStyle":"3D"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "brandtoon_session_id", Value: "session-v7"})
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	if persistedConfig.ArtisticStyle != avatarconfigdomain.ArtisticStyle3D {
		t.Fatalf("expected persisted 3D style, got %s", persistedConfig.ArtisticStyle)
	}

	var payload struct {
		AvatarConfig struct {
			AvatarID      string `json:"avatarId"`
			ArtisticStyle string `json:"artisticStyle"`
			Prompt        string `json:"prompt"`
		} `json:"avatar_config"`
	}
	decodeAvatarConfigResponse(t, recorder, &payload)

	if payload.AvatarConfig.AvatarID != "avatar-v7" {
		t.Fatalf("expected avatar-v7, got %s", payload.AvatarConfig.AvatarID)
	}
}

func TestAvatarConfigPutRejectsPromptLongerThanTwoHundredFiftySixCharacters(t *testing.T) {
	t.Parallel()

	server := newAuthenticatedAvatarConfigTestServer(t, avatarconfighttp.RouteDependencies{
		AvatarConfigRepo: &avatarconfigmocks.AvatarConfigRepositoryMock{},
		AvatarRepo: &avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
		},
	})

	request := httptest.NewRequest(
		http.MethodPut,
		"/creative-studio/avatar_configs/avatar-v7",
		bytes.NewBufferString(`{"prompt":"`+strings.Repeat("a", 257)+`","artisticStyle":"2D"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "brandtoon_session_id", Value: "session-v7"})
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", recorder.Code)
	}
}

func newAuthenticatedAvatarConfigTestServer(
	t *testing.T,
	deps avatarconfighttp.RouteDependencies,
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

	return newAvatarConfigTestServer(deps)
}

func newAvatarConfigTestServer(deps avatarconfighttp.RouteDependencies) http.Handler {
	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig("Test API", "1.0.0"))
	avatarconfighttp.RegisterRoutes(api, router, deps)
	return router
}

func decodeAvatarConfigResponse(
	t *testing.T,
	recorder *httptest.ResponseRecorder,
	target any,
) {
	t.Helper()

	if err := json.Unmarshal(recorder.Body.Bytes(), target); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}
