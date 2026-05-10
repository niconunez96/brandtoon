package test

import (
	"context"
	"testing"

	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarmocks "brandtoonapi/bounded_contexts/creative_studio/avatar/domain/mocks"
	avatarusecases "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases"
	avatardto "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases/dto"
)

func TestGetAvatarReturnsOwnedAvatarWithoutEmbeddedOptions(t *testing.T) {
	t.Parallel()

	avatar, err := avatarusecases.GetAvatar(
		context.Background(),
		avatarusecases.GetAvatarQuery{AvatarID: "avatar-v7", UserID: "user-v7"},
		&avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				storedAvatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &storedAvatar, nil
			},
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	avatarDTO := avatardto.AvatarDTO(avatar)

	if avatarDTO.ID != "avatar-v7" {
		t.Fatalf("expected avatar-v7, got %s", avatarDTO.ID)
	}

	if avatarDTO.Name != "Studio Hero" {
		t.Fatalf("expected Studio Hero, got %s", avatarDTO.Name)
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
