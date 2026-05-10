package test

import (
	"context"
	"errors"
	"testing"

	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarmocks "brandtoonapi/bounded_contexts/creative_studio/avatar/domain/mocks"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	avatarconfigmocks "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain/mocks"
	avatarconfigusecases "brandtoonapi/bounded_contexts/creative_studio/avatar_config/useCases"
	avatarconfigdto "brandtoonapi/bounded_contexts/creative_studio/avatar_config/useCases/dto"
)

func TestGetAvatarConfigReturnsNilWhenDraftDoesNotExistYet(t *testing.T) {
	t.Parallel()

	result, err := avatarconfigusecases.GetAvatarConfig(
		context.Background(),
		avatarconfigusecases.GetAvatarConfigQuery{
			AvatarID: "avatar-v7",
			UserID:   "user-v7",
		},
		&avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
		},
		&avatarconfigmocks.AvatarConfigRepositoryMock{},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if result != nil {
		t.Fatalf("expected nil avatar config, got %+v", result)
	}
}

func TestGetAvatarConfigReturnsStoredDraftForOwnedAvatar(t *testing.T) {
	t.Parallel()

	result, err := avatarconfigusecases.GetAvatarConfig(
		context.Background(),
		avatarconfigusecases.GetAvatarConfigQuery{
			AvatarID: "avatar-v7",
			UserID:   "user-v7",
		},
		&avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
		},
		&avatarconfigmocks.AvatarConfigRepositoryMock{
			FindByAvatarIDFunc: func(ctx context.Context, avatarID string) (*avatarconfigdomain.AvatarConfig, error) {
				config := avatarconfigdomain.NewAvatarConfig(
					avatarID,
					"Energetic mascot with bold shapes",
					avatarconfigdomain.ArtisticStyle3D,
					avatarconfigdomain.PersonalityBold,
				)
				return &config, nil
			},
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	configDTO := (*avatarconfigdto.AvatarConfigDTO)(result)

	if configDTO == nil || configDTO.ArtisticStyle != string(avatarconfigdomain.ArtisticStyle3D) {
		t.Fatalf("expected stored 3D config, got %+v", configDTO)
	}

	if configDTO.Personality != string(avatarconfigdomain.PersonalityBold) {
		t.Fatalf("expected stored Bold personality, got %+v", configDTO)
	}
}

func TestGetAvatarConfigReturnsAvatarNotFoundWhenAvatarIsNotOwned(t *testing.T) {
	t.Parallel()

	_, err := avatarconfigusecases.GetAvatarConfig(
		context.Background(),
		avatarconfigusecases.GetAvatarConfigQuery{
			AvatarID: "avatar-v7",
			UserID:   "user-v7",
		},
		&avatarmocks.AvatarRepositoryMock{},
		&avatarconfigmocks.AvatarConfigRepositoryMock{},
	)
	if !errors.Is(err, avatarconfigusecases.ErrAvatarNotFound) {
		t.Fatalf("expected ErrAvatarNotFound, got %v", err)
	}
}
