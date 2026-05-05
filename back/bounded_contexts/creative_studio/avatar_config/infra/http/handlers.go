package avatarconfighttp

import (
	avatarhttp "brandtoonapi/bounded_contexts/creative_studio/avatar/infra/http"
	avatarusecases "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	avatarconfigusecases "brandtoonapi/bounded_contexts/creative_studio/avatar_config/useCases"
	"context"
	"errors"
	stdhttp "net/http"

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
			response.Body.AvatarConfig = serializeAvatarConfigDTO(*avatarConfig)
		}

		return response, nil
	}
}

func buildGenerateAvatarOptionsHandler(
	deps RouteDependencies,
) stdhttp.HandlerFunc {
	return func(writer stdhttp.ResponseWriter, request *stdhttp.Request) {
		userMetadata, err := avatarhttp.RequireAuthUserMetadata(request.Context())
		if err != nil {
			writer.WriteHeader(stdhttp.StatusUnauthorized)
			return
		}

		avatarID := request.PathValue("avatarId")
		err = avatarusecases.GenerateAvatarOptions(
			request.Context(),
			avatarusecases.GenerateAvatarOptionsCommand{
				AvatarID: avatarID,
				UserID:   userMetadata.UserId,
			},
			avatarusecases.GenerateAvatarOptionsDependencies{
				AvatarConfigRepo: deps.AvatarConfigRepo,
				AvatarRepo:       deps.AvatarRepo,
				EventBus:         deps.EventBus,
			},
		)
		if err != nil {
			switch {
			case errors.Is(err, avatarusecases.ErrAvatarNotFound),
				errors.Is(err, avatarusecases.ErrAvatarConfigNotFound):
				writer.WriteHeader(stdhttp.StatusNotFound)
			case errors.Is(err, avatarconfigdomain.ErrInvalidArtisticStyle),
				errors.Is(err, avatarconfigdomain.ErrInvalidPersonality):
				writer.WriteHeader(stdhttp.StatusUnprocessableEntity)
			default:
				writer.WriteHeader(stdhttp.StatusInternalServerError)
			}
			return
		}

		writer.WriteHeader(stdhttp.StatusOK)
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
				Personality:   input.Body.Personality,
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
		response.Body.AvatarConfig = serializeAvatarConfigDTO(avatarConfig)
		return response, nil
	}
}

func mapAvatarConfigError(err error) error {
	if errors.Is(err, avatarconfigusecases.ErrAvatarNotFound) {
		return huma.Error404NotFound("avatar not found")
	}

	if errors.Is(err, avatarconfigusecases.ErrAvatarConfigNotFound) {
		return huma.Error404NotFound("avatar config not found")
	}

	if errors.Is(err, avatarconfigdomain.ErrInvalidArtisticStyle) {
		return huma.Error422UnprocessableEntity("invalid artistic style")
	}

	if errors.Is(err, avatarconfigdomain.ErrInvalidPersonality) {
		return huma.Error422UnprocessableEntity("invalid personality")
	}

	return err
}
