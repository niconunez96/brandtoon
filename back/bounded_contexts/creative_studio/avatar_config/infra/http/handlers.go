package avatarconfighttp

import (
	avatarhttp "brandtoonapi/bounded_contexts/creative_studio/avatar/infra/http"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	avatarconfigusecases "brandtoonapi/bounded_contexts/creative_studio/avatar_config/useCases"
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
)

func buildGetAvatarConfigHandler(
	deps RouteDependencies,
) func(ctx context.Context, input *avatarConfigPathInput) (*avatarConfigOutput, error) {
	return func(ctx context.Context, input *avatarConfigPathInput) (*avatarConfigOutput, error) {
		userMetadata, err := avatarhttp.RequireAuthUserMetadata(ctx)
		if err != nil {
			return nil, err
		}

		avatarConfig, err := avatarconfigusecases.GetAvatarConfig(
			ctx,
			avatarconfigusecases.GetAvatarConfigQuery{
				AvatarID: input.AvatarID,
				UserID:   userMetadata.UserId,
			},
			deps.AvatarRepo,
			deps.AvatarConfigRepo,
		)
		if err != nil {
			return nil, mapAvatarConfigError(err)
		}

		response := &avatarConfigOutput{}
		if avatarConfig != nil {
			response.Body.AvatarConfig = toAvatarConfigPayload(*avatarConfig)
		}

		return response, nil
	}
}

func buildUpdateAvatarConfigHandler(
	deps RouteDependencies,
) func(ctx context.Context, input *updateAvatarConfigInput) (*avatarConfigOutput, error) {
	return func(ctx context.Context, input *updateAvatarConfigInput) (*avatarConfigOutput, error) {
		userMetadata, err := avatarhttp.RequireAuthUserMetadata(ctx)
		if err != nil {
			return nil, err
		}

		avatarConfig, err := avatarconfigusecases.UpdateAvatarConfig(
			ctx,
			avatarconfigusecases.UpdateAvatarConfigCommand{
				AvatarID:      input.AvatarID,
				ArtisticStyle: input.Body.ArtisticStyle,
				Prompt:        input.Body.Prompt,
				UserID:        userMetadata.UserId,
			},
			deps.AvatarRepo,
			deps.AvatarConfigRepo,
		)
		if err != nil {
			return nil, mapAvatarConfigError(err)
		}

		response := &avatarConfigOutput{}
		response.Body.AvatarConfig = toAvatarConfigPayload(avatarConfig)
		return response, nil
	}
}

func mapAvatarConfigError(err error) error {
	if errors.Is(err, avatarconfigusecases.ErrAvatarNotFound) {
		return huma.Error404NotFound("avatar not found")
	}

	if errors.Is(err, avatarconfigdomain.ErrInvalidArtisticStyle) {
		return huma.Error422UnprocessableEntity("invalid artistic style")
	}

	return err
}
