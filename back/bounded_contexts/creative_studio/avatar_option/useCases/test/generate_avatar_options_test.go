package test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarmocks "brandtoonapi/bounded_contexts/creative_studio/avatar/domain/mocks"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	avatarconfigmocks "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain/mocks"
	avataroptiondomain "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain"
	avataroptionmocks "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain/mocks"
	avataroptionusecases "brandtoonapi/bounded_contexts/creative_studio/avatar_option/useCases"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	sharedmocks "brandtoonapi/bounded_contexts/shared/domain/mocks"

	"github.com/google/uuid"
)

func TestGenerateAvatarOptionsAcknowledgesAfterCreatingPendingRowsThenMarksDone(t *testing.T) {
	t.Parallel()

	createdPending := []avataroptiondomain.AvatarOption{}
	markedDone := make(chan string, 2)
	eventBus := &sharedmocks.EventBusMock{}

	err := avataroptionusecases.GenerateAvatarOptions(
		context.Background(),
		avataroptionusecases.GenerateAvatarOptionsCommand{AvatarID: "avatar-v7", UserID: "user-v7"},
		avataroptionusecases.GenerateAvatarOptionsDependencies{
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
			AvatarGenerator: &avataroptionmocks.AvatarGeneratorMock{
				GenerateOptionsFunc: func(ctx context.Context, prompt string, count int) ([]avataroptiondomain.GeneratedAvatarImage, error) {
					for _, fragment := range []string{
						"Studio Hero",
						"Energetic mascot",
						"2D",
						"Follow the artistic style STRICTLY with no style drift.",
						"Produce exactly one isolated single centered character",
						"No background. No environment. No floor. No props. No text.",
						"Avoid any 3D, CGI, clay, or realistic lighting treatment.",
					} {
						if !strings.Contains(prompt, fragment) {
							t.Fatalf("expected prompt to contain %q, got %q", fragment, prompt)
						}
					}
					if strings.Contains(prompt, "Friendly") {
						t.Fatalf("expected prompt to omit personality, got %q", prompt)
					}
					if count != 1 {
						t.Fatalf("expected count 1, got %d", count)
					}

					return []avataroptiondomain.GeneratedAvatarImage{
						{ContentType: "image/png", Data: []byte("image-one")},
					}, nil
				},
			},
			AvatarRepo: &avatarmocks.AvatarRepositoryMock{
				FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
					avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
					return &avatar, nil
				},
			},
			AvatarOptionRepo: &avataroptionmocks.AvatarOptionRepositoryMock{
				CreatePendingBatchFunc: func(ctx context.Context, options []avataroptiondomain.AvatarOption) error {
					createdPending = append([]avataroptiondomain.AvatarOption(nil), options...)
					return nil
				},
				MarkDoneFunc: func(ctx context.Context, id string, href string) error {
					markedDone <- id
					if !strings.HasPrefix(href, "http://127.0.0.1:8888/files/avatars/avatar-v7/options/") {
						t.Fatalf("unexpected href %q", href)
					}
					return nil
				},
				MarkFailedFunc: func(ctx context.Context, id string) error {
					t.Fatalf("did not expect mark failed for %s", id)
					return nil
				},
			},
			EventBus: eventBus,
			FileStorage: &sharedmocks.FileStorageMock{
				StoreFunc: func(ctx context.Context, input shareddomain.StoreFileInput) (shareddomain.StoredFile, error) {
					return shareddomain.StoredFile{
						PublicURL: "http://127.0.0.1:8888/files/" + input.Directory + "/" + input.Name + ".png",
					}, nil
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(createdPending) != 1 {
		t.Fatalf("expected 1 pending row, got %d", len(createdPending))
	}

	for _, option := range createdPending {
		parsedID, parseErr := uuid.Parse(option.ID)
		if parseErr != nil || parsedID.Version() != 7 {
			t.Fatalf("expected UUID v7 option ID, got %q (%v)", option.ID, parseErr)
		}
		if option.Status != avataroptiondomain.StatusPending {
			t.Fatalf("expected pending row, got %+v", option)
		}
	}

	for range 1 {
		select {
		case <-markedDone:
		case <-time.After(2 * time.Second):
			t.Fatalf("expected done transition")
		}
	}

	if len(eventBus.PublishedEvents) != 1 {
		t.Fatalf("expected one completion event, got %d", len(eventBus.PublishedEvents))
	}

	completedEvent, ok := eventBus.PublishedEvents[0].(avataroptionusecases.AvatarGenerationCompletedEvent)
	if !ok {
		t.Fatalf("expected AvatarGenerationCompletedEvent, got %T", eventBus.PublishedEvents[0])
	}

	if completedEvent.OutcomeValue != "SUCCESS" {
		t.Fatalf("expected success outcome, got %q", completedEvent.OutcomeValue)
	}
}

func TestGenerateAvatarOptionsMarksPendingRowsFailedWhenGenerationFails(t *testing.T) {
	t.Parallel()

	failedIDs := make(chan string, 2)
	eventBus := &sharedmocks.EventBusMock{}
	err := avataroptionusecases.GenerateAvatarOptions(
		context.Background(),
		avataroptionusecases.GenerateAvatarOptionsCommand{AvatarID: "avatar-v7", UserID: "user-v7"},
		avataroptionusecases.GenerateAvatarOptionsDependencies{
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
			AvatarGenerator: &avataroptionmocks.AvatarGeneratorMock{
				GenerateOptionsFunc: func(ctx context.Context, prompt string, count int) ([]avataroptiondomain.GeneratedAvatarImage, error) {
					return nil, errors.New("boom")
				},
			},
			AvatarRepo: &avatarmocks.AvatarRepositoryMock{
				FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
					avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
					return &avatar, nil
				},
			},
			AvatarOptionRepo: &avataroptionmocks.AvatarOptionRepositoryMock{
				CreatePendingBatchFunc: func(ctx context.Context, options []avataroptiondomain.AvatarOption) error {
					return nil
				},
				MarkFailedFunc: func(ctx context.Context, id string) error {
					failedIDs <- id
					return nil
				},
			},
			EventBus:    eventBus,
			FileStorage: &sharedmocks.FileStorageMock{},
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	for range 1 {
		select {
		case <-failedIDs:
		case <-time.After(2 * time.Second):
			t.Fatalf("expected failed transition")
		}
	}

	if len(eventBus.PublishedEvents) != 1 {
		t.Fatalf("expected one completion event, got %d", len(eventBus.PublishedEvents))
	}

	completedEvent, ok := eventBus.PublishedEvents[0].(avataroptionusecases.AvatarGenerationCompletedEvent)
	if !ok {
		t.Fatalf("expected AvatarGenerationCompletedEvent, got %T", eventBus.PublishedEvents[0])
	}

	if completedEvent.OutcomeValue != "FAILURE" {
		t.Fatalf("expected failure outcome, got %q", completedEvent.OutcomeValue)
	}
}

func TestGenerateAvatarOptionsUses3DAntiPatternPrompt(t *testing.T) {
	t.Parallel()

	promptChecked := make(chan struct{}, 1)
	err := avataroptionusecases.GenerateAvatarOptions(
		context.Background(),
		avataroptionusecases.GenerateAvatarOptionsCommand{AvatarID: "avatar-v7", UserID: "user-v7"},
		avataroptionusecases.GenerateAvatarOptionsDependencies{
			AvatarConfigRepo: &avatarconfigmocks.AvatarConfigRepositoryMock{
				FindByAvatarIDFunc: func(ctx context.Context, avatarID string) (*avatarconfigdomain.AvatarConfig, error) {
					config := avatarconfigdomain.NewAvatarConfig(
						avatarID,
						"Space captain mascot",
						avatarconfigdomain.ArtisticStyle3D,
						avatarconfigdomain.PersonalityFriendly,
					)
					return &config, nil
				},
			},
			AvatarGenerator: &avataroptionmocks.AvatarGeneratorMock{
				GenerateOptionsFunc: func(ctx context.Context, prompt string, count int) ([]avataroptiondomain.GeneratedAvatarImage, error) {
					if !strings.Contains(prompt, "Avoid flat 2D, vector, or inked outline rendering.") {
						t.Fatalf("expected 3D anti-patterns in prompt, got %q", prompt)
					}
					if strings.Contains(prompt, "Friendly") {
						t.Fatalf("expected prompt to omit personality, got %q", prompt)
					}
					promptChecked <- struct{}{}
					return []avataroptiondomain.GeneratedAvatarImage{
						{ContentType: "image/png", Data: []byte("image-one")},
					}, nil
				},
			},
			AvatarRepo: &avatarmocks.AvatarRepositoryMock{
				FindOwnedByIDFunc: func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error) {
					avatar := avatardomain.NewAvatar(avatarID, userID, "Studio Hero")
					return &avatar, nil
				},
			},
			AvatarOptionRepo: &avataroptionmocks.AvatarOptionRepositoryMock{
				CreatePendingBatchFunc: func(ctx context.Context, options []avataroptiondomain.AvatarOption) error { return nil },
				MarkDoneFunc:           func(ctx context.Context, id string, href string) error { return nil },
				MarkFailedFunc:         func(ctx context.Context, id string) error { return nil },
			},
			FileStorage: &sharedmocks.FileStorageMock{
				StoreFunc: func(ctx context.Context, input shareddomain.StoreFileInput) (shareddomain.StoredFile, error) {
					return shareddomain.StoredFile{
						PublicURL: "http://127.0.0.1:8888/files/" + input.Directory + "/" + input.Name + ".png",
					}, nil
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	select {
	case <-promptChecked:
	case <-time.After(2 * time.Second):
		t.Fatal("expected prompt check to run")
	}
}
