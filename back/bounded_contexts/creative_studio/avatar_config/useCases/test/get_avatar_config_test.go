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
				avatar := avatardomain.NewAvatarWithOptions(avatarID, userID, "Studio Hero", nil)
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
				avatar := avatardomain.NewAvatarWithOptions(
					avatarID,
					userID,
					"Studio Hero",
					[]avatardomain.AvatarOption{{ID: "option-v7", Href: "https://cdn.brandtoon.local/options/1.png", Selected: false}},
				)
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

	if result == nil || result.ArtisticStyle != string(avatarconfigdomain.ArtisticStyle3D) {
		t.Fatalf("expected stored 3D config, got %+v", result)
	}

	if result.Personality != string(avatarconfigdomain.PersonalityBold) {
		t.Fatalf("expected stored Bold personality, got %+v", result)
	}

	if len(result.AvatarOptions) != 1 {
		t.Fatalf("expected 1 avatar option, got %d", len(result.AvatarOptions))
	}

	if result.AvatarOptions[0].ID != "option-v7" {
		t.Fatalf("expected avatar option id to be preserved, got %+v", result.AvatarOptions[0])
	}

	if !result.AvatarOptions[0].Selected {
		t.Fatalf("expected read path to normalize the first option as selected, got %+v", result.AvatarOptions[0])
	}
}

func TestGetAvatarConfigNormalizesMultipleSelectedOptionsForReads(t *testing.T) {
	t.Parallel()

	result, err := avatarconfigusecases.GetAvatarConfig(
		context.Background(),
		avatarconfigusecases.GetAvatarConfigQuery{
			AvatarID: "avatar-v7",
			UserID:   "user-v7",
		},
		&avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatarWithOptions(
					avatarID,
					userID,
					"Studio Hero",
					[]avatardomain.AvatarOption{
						{ID: "option-1", Href: "https://cdn.brandtoon.local/options/1.png", Selected: false},
						{ID: "option-2", Href: "https://cdn.brandtoon.local/options/2.png", Selected: true},
						{ID: "option-3", Href: "https://cdn.brandtoon.local/options/3.png", Selected: true},
					},
				)
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

	if len(result.AvatarOptions) != 3 {
		t.Fatalf("expected 3 avatar options, got %d", len(result.AvatarOptions))
	}

	selectedCount := 0
	for _, option := range result.AvatarOptions {
		if option.Selected {
			selectedCount++
		}
	}

	if selectedCount != 1 {
		t.Fatalf("expected exactly one selected option after normalization, got %d with %+v", selectedCount, result.AvatarOptions)
	}

	if !result.AvatarOptions[1].Selected {
		t.Fatalf("expected first pre-selected option to remain selected, got %+v", result.AvatarOptions)
	}

	if result.AvatarOptions[2].Selected {
		t.Fatalf("expected later selected option to be cleared, got %+v", result.AvatarOptions)
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
