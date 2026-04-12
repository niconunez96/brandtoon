package test

import (
	"context"
	"testing"

	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarmocks "brandtoonapi/bounded_contexts/creative_studio/avatar/domain/mocks"
	avatarusecases "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases"
)

func TestCreateAvatarGeneratesIDNormalizesNameAndPersistsAvatar(t *testing.T) {
	t.Parallel()

	createdAvatar := avatardomain.Avatar{}
	result, err := avatarusecases.CreateAvatar(
		context.Background(),
		avatarusecases.CreateAvatarCommand{
			Name:   "  Nico Hero  ",
			UserID: "user-v7",
		},
		&avatarmocks.AvatarRepositoryMock{
			CreateFunc: func(ctx context.Context, avatar avatardomain.Avatar) error {
				createdAvatar = avatar
				return nil
			},
		},
		func() (string, error) { return "avatar-v7", nil },
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if result.ID != "avatar-v7" {
		t.Fatalf("expected avatar-v7, got %s", result.ID)
	}

	if result.Name != "Nico Hero" {
		t.Fatalf("expected normalized name, got %q", result.Name)
	}

	if createdAvatar.UserID != "user-v7" {
		t.Fatalf("expected user-v7, got %s", createdAvatar.UserID)
	}
}
