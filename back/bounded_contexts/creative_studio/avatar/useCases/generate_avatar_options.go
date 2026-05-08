package avatarusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	"context"
	"fmt"
	"time"
)

const AvatarGenerationCompletedEventName = "avatar-generation.completed"
const generateAvatarOptionsJobTimeout = 30 * time.Second

type GenerateAvatarOptionsCommand struct {
	AvatarID string
	UserID   string
}

type AvatarGenerationCompletedEvent struct {
	AvatarIDValue   string `json:"avatarId"`
	AvatarNameValue string `json:"avatarName"`
	UserIDValue     string `json:"userId"`
}

func (e AvatarGenerationCompletedEvent) EventName() string {
	return AvatarGenerationCompletedEventName
}

func (e AvatarGenerationCompletedEvent) UserID() string {
	return e.UserIDValue
}

type GenerateAvatarOptionsDependencies struct {
	AvatarConfigRepo avatarconfigdomain.AvatarConfigRepository
	AvatarRepo       avatardomain.AvatarRepository
	EventBus         shareddomain.EventBus
}

func GenerateAvatarOptions(
	ctx context.Context,
	cmd GenerateAvatarOptionsCommand,
	deps GenerateAvatarOptionsDependencies,
) error {
	avatar, err := deps.AvatarRepo.FindOwnedByID(ctx, cmd.AvatarID, cmd.UserID)
	if err != nil {
		return err
	}

	if avatar == nil {
		return ErrAvatarNotFound
	}

	avatarConfig, err := deps.AvatarConfigRepo.FindByAvatarID(ctx, cmd.AvatarID)
	if err != nil {
		return err
	}

	if avatarConfig == nil {
		return ErrAvatarConfigNotFound
	}

	go func(currentAvatar avatardomain.Avatar) {
		jobCtx, cancel := context.WithTimeout(context.Background(), generateAvatarOptionsJobTimeout)
		defer cancel()

		generatedOptions, err := buildGeneratedAvatarOptions(cmd.AvatarID)
		if err != nil {
			return
		}

		nextAvatar := currentAvatar.AddOptions(generatedOptions)
		if err := deps.AvatarRepo.UpdateOptions(jobCtx, cmd.AvatarID, cmd.UserID, nextAvatar.AvatarOptions); err != nil {
			return
		}

		if deps.EventBus == nil {
			return
		}

		deps.EventBus.Publish(AvatarGenerationCompletedEvent{
			AvatarIDValue:   cmd.AvatarID,
			AvatarNameValue: currentAvatar.Name,
			UserIDValue:     cmd.UserID,
		})
	}(*avatar)

	return nil
}

func buildGeneratedAvatarOptions(avatarID string) ([]avatardomain.AvatarOption, error) {
	options := make([]avatardomain.AvatarOption, 0, 4)
	for index := range 2 {
		optionID, err := shareddomain.GenerateUUIDv7()
		if err != nil {
			return nil, err
		}

		options = append(options, avatardomain.AvatarOption{
			ID:       optionID,
			Href:     fmt.Sprintf("https://cdn.brandtoon.local/avatars/%s/options/%d.png", avatarID, index+1),
			Selected: false,
		})
	}

	return options, nil
}
