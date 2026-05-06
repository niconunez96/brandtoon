package test

import (
	"context"
	"testing"

	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarmocks "brandtoonapi/bounded_contexts/creative_studio/avatar/domain/mocks"
	avatarusecases "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases"
)

func TestGetAvatarReturnsOwnedAvatarWithNormalizedOptions(t *testing.T) {
	t.Parallel()

	avatar, err := avatarusecases.GetAvatar(
		context.Background(),
		avatarusecases.GetAvatarQuery{AvatarID: "avatar-v7", UserID: "user-v7"},
		&avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				storedAvatar := avatardomain.NewAvatarWithOptions(
					avatarID,
					userID,
					"Studio Hero",
					[]avatardomain.AvatarOption{
						{ID: "option-1", Href: "https://cdn.brandtoon.local/avatars/avatar-v7/options/1.png", Selected: false},
						{ID: "option-2", Href: "https://cdn.brandtoon.local/avatars/avatar-v7/options/2.png", Selected: true},
					},
				)
				return &storedAvatar, nil
			},
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if avatar.ID != "avatar-v7" {
		t.Fatalf("expected avatar-v7, got %s", avatar.ID)
	}

	if avatar.Name != "Studio Hero" {
		t.Fatalf("expected Studio Hero, got %s", avatar.Name)
	}

	if len(avatar.AvatarOptions) != 2 {
		t.Fatalf("expected 2 avatar options, got %d", len(avatar.AvatarOptions))
	}

	if !avatar.AvatarOptions[1].Selected {
		t.Fatalf("expected selected option to remain selected, got %+v", avatar.AvatarOptions[1])
	}
}

func TestGetAvatarReturnsNotFoundWhenAvatarDoesNotExist(t *testing.T) {
	t.Parallel()

	_, err := avatarusecases.GetAvatar(
		context.Background(),
		avatarusecases.GetAvatarQuery{AvatarID: "avatar-v7", UserID: "user-v7"},
		&avatarmocks.AvatarRepositoryMock{},
	)
	if err != avatarusecases.ErrAvatarNotFound {
		t.Fatalf("expected ErrAvatarNotFound, got %v", err)
	}
}
