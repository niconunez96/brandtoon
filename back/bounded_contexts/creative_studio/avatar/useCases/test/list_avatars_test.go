package test

import (
	"context"
	"testing"

	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarmocks "brandtoonapi/bounded_contexts/creative_studio/avatar/domain/mocks"
	avatarusecases "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases"
)

func TestListAvatarsReturnsPerUserListing(t *testing.T) {
	t.Parallel()

	result, err := avatarusecases.ListAvatars(
		context.Background(),
		avatarusecases.ListAvatarsQuery{UserID: "user-v7"},
		&avatarmocks.AvatarRepositoryMock{
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
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 avatars, got %d", len(result))
	}
}
