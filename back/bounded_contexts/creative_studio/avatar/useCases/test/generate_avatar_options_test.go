package test

import (
	"context"
	"errors"
	"testing"
	"time"

	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarmocks "brandtoonapi/bounded_contexts/creative_studio/avatar/domain/mocks"
	avatarusecases "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	avatarconfigmocks "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain/mocks"
	sharedmocks "brandtoonapi/bounded_contexts/shared/domain/mocks"

	"github.com/google/uuid"
)

func TestGenerateAvatarOptionsAcknowledgesThenPersistsAndPublishes(t *testing.T) {
	t.Parallel()

	persistedOptions := []avatardomain.AvatarOption{}
	persistedSignal := make(chan struct{}, 1)
	eventBus := &sharedmocks.EventBusMock{}

	err := avatarusecases.GenerateAvatarOptions(
		context.Background(),
		avatarusecases.GenerateAvatarOptionsCommand{AvatarID: "avatar-v7", UserID: "user-v7"},
		avatarusecases.GenerateAvatarOptionsDependencies{
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
					avatar := avatardomain.NewAvatarWithOptions(
						avatarID,
						userID,
						"Studio Hero",
						[]avatardomain.AvatarOption{
							{Href: "https://cdn.brandtoon.local/avatars/avatar-v7/options/legacy.png", Selected: false},
						},
					)
					return &avatar, nil
				},
				UpdateOptionsFunc: func(ctx context.Context, avatarID string, userID string, options []avatardomain.AvatarOption) error {
					persistedOptions = options
					if len(eventBus.PublishedEvents) != 0 {
						t.Fatalf("expected publish after persist")
					}
					persistedSignal <- struct{}{}
					return nil
				},
			},
			EventBus: eventBus,
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	select {
	case <-persistedSignal:
	case <-time.After(2 * time.Second):
		t.Fatalf("expected background persistence signal")
	}

	if len(persistedOptions) != 3 {
		t.Fatalf("expected existing option + 3 generated options, got %d", len(persistedOptions))
	}

	generatedIDs := make(map[string]struct{}, len(persistedOptions)-1)
	for _, option := range persistedOptions[1:] {
		parsedID, parseErr := uuid.Parse(option.ID)
		if parseErr != nil {
			t.Fatalf("expected generated option ID to be a UUID, got %q (%v)", option.ID, parseErr)
		}

		if parsedID.Version() != 7 {
			t.Fatalf("expected generated option ID to be UUID v7, got version %d", parsedID.Version())
		}

		if _, exists := generatedIDs[option.ID]; exists {
			t.Fatalf("expected generated option IDs to be unique, duplicated %q", option.ID)
		}
		generatedIDs[option.ID] = struct{}{}

		if option.Selected {
			t.Fatalf("expected generated options selected=false by default")
		}
	}

	if len(eventBus.PublishedEvents) != 1 {
		t.Fatalf("expected exactly one event, got %d", len(eventBus.PublishedEvents))
	}
}

func TestGenerateAvatarOptionsStopsWhenDraftIsMissing(t *testing.T) {
	t.Parallel()

	eventBus := &sharedmocks.EventBusMock{}
	avatarConfigRepo := &avatarconfigmocks.AvatarConfigRepositoryMock{}
	avatarRepo := &avatarmocks.AvatarRepositoryMock{
		FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
			avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
			return &avatar, nil
		},
		UpdateOptionsFunc: func(ctx context.Context, avatarID string, userID string, options []avatardomain.AvatarOption) error {
			t.Fatalf("did not expect options update when draft is missing")
			return nil
		},
	}

	err := avatarusecases.GenerateAvatarOptions(
		context.Background(),
		avatarusecases.GenerateAvatarOptionsCommand{AvatarID: "avatar-v7", UserID: "user-v7"},
		avatarusecases.GenerateAvatarOptionsDependencies{
			AvatarConfigRepo: avatarConfigRepo,
			AvatarRepo:       avatarRepo,
			EventBus:         eventBus,
		},
	)
	if !errors.Is(err, avatarusecases.ErrAvatarConfigNotFound) {
		t.Fatalf("expected ErrAvatarConfigNotFound, got %v", err)
	}

	if len(eventBus.PublishedEvents) != 0 {
		t.Fatalf("expected no events when draft is missing")
	}
}

func TestGenerateAvatarOptionsSuppressesPublishWhenPersistenceFails(t *testing.T) {
	t.Parallel()

	eventBus := &sharedmocks.EventBusMock{}
	persistAttempted := make(chan struct{}, 1)
	err := avatarusecases.GenerateAvatarOptions(
		context.Background(),
		avatarusecases.GenerateAvatarOptionsCommand{AvatarID: "avatar-v7", UserID: "user-v7"},
		avatarusecases.GenerateAvatarOptionsDependencies{
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
					persistAttempted <- struct{}{}
					return errors.New("boom")
				},
			},
			EventBus: eventBus,
		},
	)
	if err != nil {
		t.Fatalf("expected ACK path to return nil error, got %v", err)
	}

	select {
	case <-persistAttempted:
	case <-time.After(2 * time.Second):
		t.Fatalf("expected persistence attempt")
	}

	if len(eventBus.PublishedEvents) != 0 {
		t.Fatalf("expected no published event when persistence fails")
	}
}
