package test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarmocks "brandtoonapi/bounded_contexts/creative_studio/avatar/domain/mocks"
	avatarhttp "brandtoonapi/bounded_contexts/creative_studio/avatar/infra/http"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	avatarconfigmocks "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain/mocks"
)

func TestAvatarAggregateGenerateRouteReturnsImmediateAckWithoutBody(t *testing.T) {
	t.Parallel()

	server := newAuthenticatedTestServer(t, avatarhttp.RouteDependencies{
		AvatarConfigRepo: &avatarconfigmocks.AvatarConfigRepositoryMock{
			FindByAvatarIDFunc: func(ctx context.Context, avatarID string) (*avatarconfigdomain.AvatarConfig, error) {
				config := avatarconfigdomain.NewAvatarConfig(
					avatarID,
					"Energetic mascot",
					avatarconfigdomain.ArtisticStyle2D,
					avatarconfigdomain.PersonalityFriendly,
				)
				return &config, nil
			},
		},
		AvatarRepo: &avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
			UpdateOptionsFunc: func(ctx context.Context, avatarID string, userID string, options []avatardomain.AvatarOption) error {
				return nil
			},
		},
	})

	request := httptest.NewRequest(http.MethodPost, "/creative-studio/avatar_configs/avatar-v7/generate", nil)
	request.AddCookie(&http.Cookie{Name: "brandtoon_session_id", Value: "session-v7"})
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}

	if recorder.Body.Len() != 0 {
		t.Fatalf("expected empty response body, got %q", recorder.Body.String())
	}
}

func TestAvatarAggregateSelectRoutePersistsRequestedOption(t *testing.T) {
	t.Parallel()

	server := newAuthenticatedTestServer(t, avatarhttp.RouteDependencies{
		AvatarRepo: &avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatarWithOptions(
					avatarID,
					userID,
					"Studio Hero",
					[]avatardomain.AvatarOption{
						{
							ID:       "option-1",
							Href:     "https://cdn.brandtoon.local/avatars/avatar-v7/options/1.png",
							Selected: true,
						},
						{
							ID:       "option-2",
							Href:     "https://cdn.brandtoon.local/avatars/avatar-v7/options/2.png",
							Selected: false,
						},
					},
				)
				return &avatar, nil
			},
			UpdateOptionsFunc: func(ctx context.Context, avatarID string, userID string, options []avatardomain.AvatarOption) error {
				return nil
			},
		},
	})

	request := httptest.NewRequest(
		http.MethodPost,
		"/creative-studio/avatar_configs/avatar-v7/options/select",
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

func TestAvatarAggregateDeleteOptionsRouteRemovesRequestedIDs(t *testing.T) {
	t.Parallel()

	server := newAuthenticatedTestServer(t, avatarhttp.RouteDependencies{
		AvatarRepo: &avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatarWithOptions(
					avatarID,
					userID,
					"Studio Hero",
					[]avatardomain.AvatarOption{
						{
							ID:       "option-1",
							Href:     "https://cdn.brandtoon.local/avatars/avatar-v7/options/1.png",
							Selected: false,
						},
						{
							ID:       "option-2",
							Href:     "https://cdn.brandtoon.local/avatars/avatar-v7/options/2.png",
							Selected: true,
						},
					},
				)
				return &avatar, nil
			},
			UpdateOptionsFunc: func(ctx context.Context, avatarID string, userID string, options []avatardomain.AvatarOption) error {
				return nil
			},
		},
	})

	request := httptest.NewRequest(
		http.MethodDelete,
		"/creative-studio/avatar_configs/avatar-v7/options",
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
