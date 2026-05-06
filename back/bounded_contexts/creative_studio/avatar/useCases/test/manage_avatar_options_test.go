package test

import (
	"context"
	"errors"
	"testing"

	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarmocks "brandtoonapi/bounded_contexts/creative_studio/avatar/domain/mocks"
	avatarusecases "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases"
)

func TestSelectAvatarOptionPersistsExactlyOneSelection(t *testing.T) {
	t.Parallel()

	persistedOptions := []avatardomain.AvatarOption{}
	result, err := avatarusecases.SelectAvatarOption(
		context.Background(),
		avatarusecases.SelectAvatarOptionCommand{
			AvatarID: "avatar-v7",
			ID:       "option-3",
			UserID:   "user-v7",
		},
		&avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatarWithOptions(
					avatarID,
					userID,
					"Studio Hero",
					[]avatardomain.AvatarOption{
						{ID: "option-1", Href: "https://cdn.brandtoon.local/options/1.png", Selected: true},
						{ID: "option-2", Href: "https://cdn.brandtoon.local/options/2.png", Selected: true},
						{ID: "option-3", Href: "https://cdn.brandtoon.local/options/3.png", Selected: false},
					},
				)
				return &avatar, nil
			},
			UpdateOptionsFunc: func(ctx context.Context, avatarID string, userID string, options []avatardomain.AvatarOption) error {
				persistedOptions = append([]avatardomain.AvatarOption(nil), options...)
				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(persistedOptions) != 3 {
		t.Fatalf("expected 3 persisted options, got %d", len(persistedOptions))
	}

	selectedCount := 0
	for _, option := range persistedOptions {
		if option.Selected {
			selectedCount++
		}
	}

	if selectedCount != 1 {
		t.Fatalf("expected exactly one persisted selected option, got %d with %+v", selectedCount, persistedOptions)
	}

	if !persistedOptions[2].Selected {
		t.Fatalf("expected requested option to become selected, got %+v", persistedOptions)
	}

	if result.AvatarID != "avatar-v7" || len(result.AvatarOptions) != 3 {
		t.Fatalf("expected updated avatar option DTO, got %+v", result)
	}
}

func TestSelectAvatarOptionReturnsNotFoundWhenOptionIDIsMissing(t *testing.T) {
	t.Parallel()

	_, err := avatarusecases.SelectAvatarOption(
		context.Background(),
		avatarusecases.SelectAvatarOptionCommand{
			AvatarID: "avatar-v7",
			ID:       "missing-option",
			UserID:   "user-v7",
		},
		&avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatarWithOptions(
					avatarID,
					userID,
					"Studio Hero",
					[]avatardomain.AvatarOption{{ID: "option-1", Href: "https://cdn.brandtoon.local/options/1.png", Selected: false}},
				)
				return &avatar, nil
			},
		},
	)
	if !errors.Is(err, avatarusecases.ErrAvatarOptionNotFound) {
		t.Fatalf("expected ErrAvatarOptionNotFound, got %v", err)
	}
}

func TestDeleteAvatarOptionsRemovesRequestedOptionsAndNormalizesSelection(t *testing.T) {
	t.Parallel()

	persistedOptions := []avatardomain.AvatarOption{}
	result, err := avatarusecases.DeleteAvatarOptions(
		context.Background(),
		avatarusecases.DeleteAvatarOptionsCommand{
			AvatarID: "avatar-v7",
			IDs:      []string{"option-1", "option-2"},
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
						{ID: "option-3", Href: "https://cdn.brandtoon.local/options/3.png", Selected: false},
					},
				)
				return &avatar, nil
			},
			UpdateOptionsFunc: func(ctx context.Context, avatarID string, userID string, options []avatardomain.AvatarOption) error {
				persistedOptions = append([]avatardomain.AvatarOption(nil), options...)
				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(persistedOptions) != 1 {
		t.Fatalf("expected 1 remaining option, got %d with %+v", len(persistedOptions), persistedOptions)
	}

	if persistedOptions[0].ID != "option-3" || !persistedOptions[0].Selected {
		t.Fatalf("expected remaining option to be normalized as selected, got %+v", persistedOptions)
	}

	if len(result.AvatarOptions) != 1 || result.AvatarOptions[0].ID != "option-3" {
		t.Fatalf("expected response DTO with remaining option, got %+v", result)
	}
}

func TestDeleteAvatarOptionsReturnsNotFoundWhenNothingMatches(t *testing.T) {
	t.Parallel()

	_, err := avatarusecases.DeleteAvatarOptions(
		context.Background(),
		avatarusecases.DeleteAvatarOptionsCommand{
			AvatarID: "avatar-v7",
			IDs:      []string{"missing-option"},
			UserID:   "user-v7",
		},
		&avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatarWithOptions(
					avatarID,
					userID,
					"Studio Hero",
					[]avatardomain.AvatarOption{{ID: "option-1", Href: "https://cdn.brandtoon.local/options/1.png", Selected: false}},
				)
				return &avatar, nil
			},
		},
	)
	if !errors.Is(err, avatarusecases.ErrAvatarOptionNotFound) {
		t.Fatalf("expected ErrAvatarOptionNotFound, got %v", err)
	}
}
