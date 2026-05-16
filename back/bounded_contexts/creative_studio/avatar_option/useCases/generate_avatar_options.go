package avataroptionusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	avataroptiondomain "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	"brandtoonapi/bounded_contexts/shared/infra/telemetry"
	"context"
	"fmt"
	"path"
	"strings"
	"time"
)

const AvatarGenerationCompletedEventName = "avatar-generation.completed"
const generateAvatarOptionsJobTimeout = 2 * time.Minute
const generatedAvatarOptionCount = 1

type GenerateAvatarOptionsCommand struct {
	AvatarID string
	UserID   string
}

type AvatarGenerationCompletedEvent struct {
	AvatarIDValue   string `json:"avatarId"`
	AvatarNameValue string `json:"avatarName"`
	OutcomeValue    string `json:"outcome"`
	UserIDValue     string `json:"userId"`
}

func (e AvatarGenerationCompletedEvent) EventName() string { return AvatarGenerationCompletedEventName }
func (e AvatarGenerationCompletedEvent) UserID() string    { return e.UserIDValue }

type GenerateAvatarOptionsDependencies struct {
	AvatarConfigRepo avatarconfigdomain.AvatarConfigRepository
	AvatarGenerator  avataroptiondomain.AvatarGenerator
	AvatarRepo       avatardomain.AvatarRepository
	AvatarOptionRepo avataroptiondomain.AvatarOptionRepository
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

	if deps.AvatarGenerator == nil || deps.FileStorage == nil || deps.AvatarOptionRepo == nil {
		return ErrAvatarGenerationUnavailable
	}

	pendingOptions, err := buildPendingAvatarOptions(cmd.AvatarID)
	if err != nil {
		return err
	}

	if err := deps.AvatarOptionRepo.CreatePendingBatch(ctx, pendingOptions); err != nil {
		return err
	}

	go processAvatarOptionsGeneration(*avatar, *avatarConfig, pendingOptions, cmd, deps)
	return nil
}

func buildPendingAvatarOptions(avatarID string) ([]avataroptiondomain.AvatarOption, error) {
	options := make([]avataroptiondomain.AvatarOption, 0, generatedAvatarOptionCount)
	for range generatedAvatarOptionCount {
		optionID, err := shareddomain.GenerateUUIDv7()
		if err != nil {
			return nil, err
		}

		option, err := avataroptiondomain.NewPendingAvatarOption(optionID, avatarID)
		if err != nil {
			return nil, err
		}

		options = append(options, option)
	}

	return options, nil
}

func processAvatarOptionsGeneration(
	avatar avatardomain.Avatar,
	avatarConfig avatarconfigdomain.AvatarConfig,
	pendingOptions []avataroptiondomain.AvatarOption,
	cmd GenerateAvatarOptionsCommand,
	deps GenerateAvatarOptionsDependencies,
) {
	jobCtx, cancel := context.WithTimeout(context.Background(), generateAvatarOptionsJobTimeout)
	defer cancel()

	generatedImages, err := deps.AvatarGenerator.GenerateOptions(
		jobCtx,
		buildAvatarGenerationPrompt(avatar, avatarConfig),
		generatedAvatarOptionCount,
	)
	if err != nil {
		telemetry.LogError("Error while generating avatar options", err)
		markAllFailed(jobCtx, pendingOptions, deps.AvatarOptionRepo)
		publishAvatarGenerationCompleted(cmd, avatar, avatarGenerationOutcomeFailure, deps.EventBus)
		return
	}

	successfullyGeneratedOptions := 0

	for index, option := range pendingOptions {
		if index >= len(generatedImages) {
			_ = deps.AvatarOptionRepo.MarkFailed(jobCtx, option.ID)
			continue
		}

		storedFile, storeErr := deps.FileStorage.Store(jobCtx, shareddomain.StoreFileInput{
			ContentType: generatedImages[index].ContentType,
			Data:        generatedImages[index].Data,
			Directory:   path.Join("avatars", avatar.ID, "options"),
			Name:        option.ID,
		})
		if storeErr != nil {
			telemetry.LogError("Error while saving avatar option in storage", storeErr)
			_ = deps.AvatarOptionRepo.MarkFailed(jobCtx, option.ID)
			continue
		}

		if err := deps.AvatarOptionRepo.MarkDone(jobCtx, option.ID, storedFile.PublicURL); err != nil {
			telemetry.LogError("Error while marking avatar option done", err)
			_ = deps.AvatarOptionRepo.MarkFailed(jobCtx, option.ID)
			continue
		}

		successfullyGeneratedOptions++
	}

	outcome := avatarGenerationOutcomeFailure
	if successfullyGeneratedOptions > 0 {
		outcome = avatarGenerationOutcomeSuccess
	}

	publishAvatarGenerationCompleted(cmd, avatar, outcome, deps.EventBus)
}

const (
	avatarGenerationOutcomeSuccess = "SUCCESS"
	avatarGenerationOutcomeFailure = "FAILURE"
)

func markAllFailed(
	ctx context.Context,
	options []avataroptiondomain.AvatarOption,
	repo avataroptiondomain.AvatarOptionRepository,
) {
	for _, option := range options {
		_ = repo.MarkFailed(ctx, option.ID)
	}
}

func publishAvatarGenerationCompleted(
	cmd GenerateAvatarOptionsCommand,
	avatar avatardomain.Avatar,
	outcome string,
	eventBus shareddomain.EventBus,
) {
	if eventBus == nil {
		return
	}

	eventBus.Publish(
		AvatarGenerationCompletedEvent{
			AvatarIDValue:   cmd.AvatarID,
			AvatarNameValue: avatar.Name,
			OutcomeValue:    outcome,
			UserIDValue:     cmd.UserID,
		},
	)
}

func buildAvatarGenerationPrompt(
	avatar avatardomain.Avatar,
	avatarConfig avatarconfigdomain.AvatarConfig,
) string {
	styleSpecificAntiPattern := "Avoid flat 2D, vector, or inked outline rendering."
	if avatarConfig.ArtisticStyle == avatarconfigdomain.ArtisticStyle2D {
		styleSpecificAntiPattern = "Avoid any 3D, CGI, clay, or realistic lighting treatment."
	}

	return strings.TrimSpace(fmt.Sprintf(
		"Create a polished brand avatar portrait for the character named %s. "+
			"Base concept: %s. Artistic style: %s. "+
			"Follow the artistic style STRICTLY with no style drift. "+
			"Produce exactly one isolated single centered character with a strong silhouette, suitable for product avatar selection. "+
			"No background. No environment. No floor. No props. No text. %s",
		avatar.Name,
		avatarConfig.Prompt,
		avatarConfig.ArtisticStyle,
		styleSpecificAntiPattern,
	))
}
