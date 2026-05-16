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

func TestUpdateAvatarConfigUpsertsDraftForOwnedAvatar(t *testing.T) {
	t.Parallel()

	persistedConfig := avatarconfigdomain.AvatarConfig{}
	result, err := avatarconfigusecases.UpdateAvatarConfig(
		context.Background(),
		avatarconfigusecases.UpdateAvatarConfigCommand{
			AvatarID:      "avatar-v7",
			ArtisticStyle: "2D",
			Personality:   "Friendly",
			Prompt:        "  ",
			UserID:        "user-v7",
		},
		&avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
		},
		&avatarconfigmocks.AvatarConfigRepositoryMock{
			UpsertFunc: func(ctx context.Context, avatarConfig avatarconfigdomain.AvatarConfig) error {
				persistedConfig = avatarConfig
				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	resultDTO := avatarconfigdto.AvatarConfigDTO(result)

	if resultDTO.Prompt != "  " {
		t.Fatalf("expected whitespace prompt to be preserved, got %q", resultDTO.Prompt)
	}

	if persistedConfig.ArtisticStyle != avatarconfigdomain.ArtisticStyle2D {
		t.Fatalf("expected 2D artistic style, got %s", persistedConfig.ArtisticStyle)
	}

	if persistedConfig.Personality != avatarconfigdomain.PersonalityFriendly {
		t.Fatalf("expected Friendly personality, got %s", persistedConfig.Personality)
	}
}

func TestUpdateAvatarConfigReturnsAvatarNotFoundWhenAvatarIsMissing(t *testing.T) {
	t.Parallel()

	_, err := avatarconfigusecases.UpdateAvatarConfig(
		context.Background(),
		avatarconfigusecases.UpdateAvatarConfigCommand{
			AvatarID:      "avatar-v7",
			ArtisticStyle: "2D",
			Personality:   "Friendly",
			Prompt:        "hello",
			UserID:        "user-v7",
		},
		&avatarmocks.AvatarRepositoryMock{},
		&avatarconfigmocks.AvatarConfigRepositoryMock{},
	)
	if !errors.Is(err, avatarconfigusecases.ErrAvatarNotFound) {
		t.Fatalf("expected ErrAvatarNotFound, got %v", err)
	}
}

func TestUpdateAvatarConfigRejectsInvalidArtisticStyle(t *testing.T) {
	t.Parallel()

	_, err := avatarconfigusecases.UpdateAvatarConfig(
		context.Background(),
		avatarconfigusecases.UpdateAvatarConfigCommand{
			AvatarID:      "avatar-v7",
			ArtisticStyle: "Clay",
			Personality:   "Friendly",
			Prompt:        "hello",
			UserID:        "user-v7",
		},
		&avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
		},
		&avatarconfigmocks.AvatarConfigRepositoryMock{},
	)
	if !errors.Is(err, avatarconfigdomain.ErrInvalidArtisticStyle) {
		t.Fatalf("expected ErrInvalidArtisticStyle, got %v", err)
	}
}

func TestUpdateAvatarConfigRejectsInvalidPersonality(t *testing.T) {
	t.Parallel()

	_, err := avatarconfigusecases.UpdateAvatarConfig(
		context.Background(),
		avatarconfigusecases.UpdateAvatarConfigCommand{
			AvatarID:      "avatar-v7",
			ArtisticStyle: "2D",
			Personality:   "Serious",
			Prompt:        "hello",
			UserID:        "user-v7",
		},
		&avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
		},
		&avatarconfigmocks.AvatarConfigRepositoryMock{},
	)
	if !errors.Is(err, avatarconfigdomain.ErrInvalidPersonality) {
		t.Fatalf("expected ErrInvalidPersonality, got %v", err)
	}
}
