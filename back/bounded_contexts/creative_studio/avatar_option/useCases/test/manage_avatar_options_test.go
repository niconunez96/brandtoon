package test

import (
	"context"
	"errors"
	"testing"

	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarmocks "brandtoonapi/bounded_contexts/creative_studio/avatar/domain/mocks"
	avataroptiondomain "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain"
	avataroptionmocks "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain/mocks"
	avataroptionusecases "brandtoonapi/bounded_contexts/creative_studio/avatar_option/useCases"
	avataroptiondto "brandtoonapi/bounded_contexts/creative_studio/avatar_option/useCases/dto"
)

func TestListAvatarOptionsReturnsOwnedAvatarOptions(t *testing.T) {
	t.Parallel()

	doneHref := "https://cdn.brandtoon.local/avatars/avatar-v7/options/1.png"
	result, err := avataroptionusecases.ListAvatarOptions(
		context.Background(),
		avataroptionusecases.ListAvatarOptionsQuery{AvatarID: "avatar-v7", UserID: "user-v7"},
		&avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
		},
		&avataroptionmocks.AvatarOptionRepositoryMock{
			ListByAvatarIDFunc: func(ctx context.Context, avatarID string) ([]avataroptiondomain.AvatarOption, error) {
				return []avataroptiondomain.AvatarOption{{
					ID:       "option-1",
					AvatarID: avatarID,
					Status:   avataroptiondomain.StatusDone,
					Selected: true,
					Href:     &doneHref,
				}}, nil
			},
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	resultDTO := avataroptiondto.AvatarOptionsDTO(result)
	if len(resultDTO.AvatarOptions) != 1 || resultDTO.AvatarOptions[0].Status != string(avataroptiondomain.StatusDone) {
		t.Fatalf("expected one done option, got %+v", resultDTO)
	}
}

func TestSelectAvatarOptionPersistsExactlyOneSelection(t *testing.T) {
	t.Parallel()

	result, err := avataroptionusecases.SelectAvatarOption(
		context.Background(),
		avataroptionusecases.SelectAvatarOptionCommand{AvatarID: "avatar-v7", ID: "option-2", UserID: "user-v7"},
		&avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
		},
		&avataroptionmocks.AvatarOptionRepositoryMock{
			SelectFunc: func(ctx context.Context, avatarID string, optionID string) error {
				if avatarID != "avatar-v7" || optionID != "option-2" {
					t.Fatalf("unexpected select args %s %s", avatarID, optionID)
				}
				return nil
			},
			ListByAvatarIDFunc: func(ctx context.Context, avatarID string) ([]avataroptiondomain.AvatarOption, error) {
				hrefOne := "https://cdn.brandtoon.local/avatars/avatar-v7/options/1.png"
				hrefTwo := "https://cdn.brandtoon.local/avatars/avatar-v7/options/2.png"
				return []avataroptiondomain.AvatarOption{
					{
						ID:       "option-1",
						AvatarID: avatarID,
						Status:   avataroptiondomain.StatusDone,
						Href:     &hrefOne,
						Selected: false,
					},
					{
						ID:       "option-2",
						AvatarID: avatarID,
						Status:   avataroptiondomain.StatusDone,
						Href:     &hrefTwo,
						Selected: true,
					},
				}, nil
			},
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(result.AvatarOptions) != 2 || !result.AvatarOptions[1].Selected {
		t.Fatalf("expected option-2 selected in response, got %+v", result)
	}
}

func TestDeleteAvatarOptionsRemovesRequestedOptions(t *testing.T) {
	t.Parallel()

	result, err := avataroptionusecases.DeleteAvatarOptions(
		context.Background(),
		avataroptionusecases.DeleteAvatarOptionsCommand{
			AvatarID: "avatar-v7",
			IDs:      []string{"option-1"},
			UserID:   "user-v7",
		},
		&avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
		},
		&avataroptionmocks.AvatarOptionRepositoryMock{
			DeleteFunc: func(ctx context.Context, avatarID string, ids []string) error {
				if len(ids) != 1 || ids[0] != "option-1" {
					t.Fatalf("unexpected delete ids %+v", ids)
				}
				return nil
			},
			ListByAvatarIDFunc: func(ctx context.Context, avatarID string) ([]avataroptiondomain.AvatarOption, error) {
				hrefTwo := "https://cdn.brandtoon.local/avatars/avatar-v7/options/2.png"
				return []avataroptiondomain.AvatarOption{
					{
						ID:       "option-2",
						AvatarID: avatarID,
						Status:   avataroptiondomain.StatusDone,
						Href:     &hrefTwo,
						Selected: true,
					},
				}, nil
			},
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(result.AvatarOptions) != 1 || result.AvatarOptions[0].ID != "option-2" {
		t.Fatalf("expected only option-2 remaining, got %+v", result)
	}
}

func TestSelectAvatarOptionReturnsNotFoundWhenRepositoryRejectsSelection(t *testing.T) {
	t.Parallel()

	_, err := avataroptionusecases.SelectAvatarOption(
		context.Background(),
		avataroptionusecases.SelectAvatarOptionCommand{AvatarID: "avatar-v7", ID: "missing-option", UserID: "user-v7"},
		&avatarmocks.AvatarRepositoryMock{
			FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
				avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
				return &avatar, nil
			},
		},
		&avataroptionmocks.AvatarOptionRepositoryMock{
			SelectFunc: func(ctx context.Context, avatarID string, optionID string) error {
				return avataroptiondomain.ErrAvatarOptionNotFound
			},
		},
	)
	if !errors.Is(err, avataroptionusecases.ErrAvatarOptionNotFound) {
		t.Fatalf("expected ErrAvatarOptionNotFound, got %v", err)
	}
}
