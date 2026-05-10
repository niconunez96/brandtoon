package avatarusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	"brandtoonapi/bounded_contexts/shared/infra/telemetry"
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"
)

const AvatarGenerationCompletedEventName = "avatar-generation.completed"
const generateAvatarOptionsJobTimeout = 2 * time.Minute
const generatedAvatarOptionCount = 2

var ErrAvatarGenerationUnavailable = errors.New("avatar generation unavailable")

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
	AvatarGenerator  avatardomain.AvatarGenerator
	AvatarRepo       avatardomain.AvatarRepository
	EventBus         shareddomain.EventBus
	FileStorage      shareddomain.FileStorage
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

	if deps.AvatarGenerator == nil || deps.FileStorage == nil {
		return ErrAvatarGenerationUnavailable
	}

	go func(currentAvatar avatardomain.Avatar, currentConfig avatarconfigdomain.AvatarConfig) {
		jobCtx, cancel := context.WithTimeout(context.Background(), generateAvatarOptionsJobTimeout)
		defer cancel()

		generatedOptions, err := buildGeneratedAvatarOptions(
			jobCtx,
			currentAvatar,
			currentConfig,
			deps.AvatarGenerator,
			deps.FileStorage,
		)
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
	}(*avatar, *avatarConfig)

	return nil
}

func buildGeneratedAvatarOptions(
	ctx context.Context,
	avatar avatardomain.Avatar,
	avatarConfig avatarconfigdomain.AvatarConfig,
	generator avatardomain.AvatarGenerator,
	fileStorage shareddomain.FileStorage,
) ([]avatardomain.AvatarOption, error) {
	prompt := buildAvatarGenerationPrompt(avatar, avatarConfig)
	generatedImages, err := generator.GenerateOptions(ctx, prompt, generatedAvatarOptionCount)
	if err != nil {
		telemetry.LogError("Error while generating images", err)
		return nil, err
	}
	telemetry.LogInfo(fmt.Sprintf("Generated Images %d", len(generatedImages)))

	options := make([]avatardomain.AvatarOption, 0, len(generatedImages))
	for _, generatedImage := range generatedImages {
		optionID, err := shareddomain.GenerateUUIDv7()
		if err != nil {
			return nil, err
		}

		storedFile, err := fileStorage.Store(ctx, shareddomain.StoreFileInput{
			ContentType: generatedImage.ContentType,
			Data:        generatedImage.Data,
			Directory:   path.Join("avatars", avatar.ID, "options"),
			Name:        optionID,
		})
		if err != nil {
			telemetry.LogError("Error while saving avatar in storage", err)
			return nil, err
		}

		options = append(options, avatardomain.AvatarOption{
			ID:       optionID,
			Href:     storedFile.PublicURL,
			Selected: false,
		})
	}

	return options, nil
}

func buildAvatarGenerationPrompt(
	avatar avatardomain.Avatar,
	avatarConfig avatarconfigdomain.AvatarConfig,
) string {
	return strings.TrimSpace(fmt.Sprintf(
		"Create a polished brand avatar portrait for the character named %s. Base concept: %s. Artistic style: %s. Personality: %s. Produce a single centered character variation with a clean background and strong silhouette, suitable for product avatar selection.",
		avatar.Name,
		avatarConfig.Prompt,
		avatarConfig.ArtisticStyle,
		avatarConfig.Personality,
	))
}
