package avataroptionhttp

import (
	avatarhttp "brandtoonapi/bounded_contexts/creative_studio/avatar/infra/http"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	avataroptiondomain "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain"
	avataroptionusecases "brandtoonapi/bounded_contexts/creative_studio/avatar_option/useCases"
	"context"
	"errors"
	stdhttp "net/http"

	"github.com/danielgtaylor/huma/v2"
)

func buildListAvatarOptionsHandler(
	deps RouteDependencies,
) func(ctx context.Context, input *listAvatarOptionsInput) (*avatarOptionsOutput, error) {
	return func(ctx context.Context, input *listAvatarOptionsInput) (*avatarOptionsOutput, error) {
		userMetadata, err := avatarhttp.RequireAuthUserMetadata(ctx)
		if err != nil {
			return nil, err
		}

		avatarOptions, err := avataroptionusecases.ListAvatarOptions(
			ctx,
			avataroptionusecases.ListAvatarOptionsQuery{AvatarID: input.AvatarID, UserID: userMetadata.UserId},
			deps.AvatarRepo,
			deps.AvatarOptionRepo,
		)
		if err != nil {
			return nil, mapAvatarOptionsError(err)
		}

		response := &avatarOptionsOutput{}
		response.Body.Avatar = avatarOptions
		return response, nil
	}
}

func buildGenerateAvatarOptionsHandler(deps RouteDependencies) stdhttp.HandlerFunc {
	return func(writer stdhttp.ResponseWriter, request *stdhttp.Request) {
		userMetadata, err := avatarhttp.RequireAuthUserMetadata(request.Context())
		if err != nil {
			writer.WriteHeader(stdhttp.StatusUnauthorized)
			return
		}

		err = avataroptionusecases.GenerateAvatarOptions(
			request.Context(),
			avataroptionusecases.GenerateAvatarOptionsCommand{
				AvatarID: request.PathValue("avatarId"),
				UserID:   userMetadata.UserId,
			},
			avataroptionusecases.GenerateAvatarOptionsDependencies{
				AvatarConfigRepo: deps.AvatarConfigRepo,
				AvatarGenerator:  deps.AvatarGenerator,
				AvatarRepo:       deps.AvatarRepo,
				AvatarOptionRepo: deps.AvatarOptionRepo,
				EventBus:         deps.EventBus,
				FileStorage:      deps.FileStorage,
			},
		)
		if err != nil {
			switch {
			case errors.Is(err, avataroptionusecases.ErrAvatarNotFound),
				errors.Is(err, avataroptionusecases.ErrAvatarConfigNotFound):
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

func buildSelectAvatarOptionHandler(
	deps RouteDependencies,
) func(ctx context.Context, input *selectAvatarOptionInput) (*avatarOptionsOutput, error) {
	return func(ctx context.Context, input *selectAvatarOptionInput) (*avatarOptionsOutput, error) {
		userMetadata, err := avatarhttp.RequireAuthUserMetadata(ctx)
		if err != nil {
			return nil, err
		}

		avatarOptions, err := avataroptionusecases.SelectAvatarOption(
			ctx,
			avataroptionusecases.SelectAvatarOptionCommand{
				AvatarID: input.AvatarID,
				ID:       input.Body.ID,
				UserID:   userMetadata.UserId,
			},
			deps.AvatarRepo,
			deps.AvatarOptionRepo,
		)
		if err != nil {
			return nil, mapAvatarOptionsError(err)
		}

		response := &avatarOptionsOutput{}
		response.Body.Avatar = avatarOptions
		return response, nil
	}
}

func buildDeleteAvatarOptionsHandler(
	deps RouteDependencies,
) func(ctx context.Context, input *deleteAvatarOptionsInput) (*avatarOptionsOutput, error) {
	return func(ctx context.Context, input *deleteAvatarOptionsInput) (*avatarOptionsOutput, error) {
		userMetadata, err := avatarhttp.RequireAuthUserMetadata(ctx)
		if err != nil {
			return nil, err
		}

		avatarOptions, err := avataroptionusecases.DeleteAvatarOptions(
			ctx,
			avataroptionusecases.DeleteAvatarOptionsCommand{
				AvatarID: input.AvatarID,
				IDs:      input.Body.IDs,
				UserID:   userMetadata.UserId,
			},
			deps.AvatarRepo,
			deps.AvatarOptionRepo,
		)
		if err != nil {
			return nil, mapAvatarOptionsError(err)
		}

		response := &avatarOptionsOutput{}
		response.Body.Avatar = avatarOptions
		return response, nil
	}
}

func mapAvatarOptionsError(err error) error {
	if errors.Is(err, avataroptionusecases.ErrAvatarNotFound) {
		return huma.Error404NotFound("avatar not found")
	}
	if errors.Is(err, avataroptionusecases.ErrAvatarOptionNotFound) ||
		errors.Is(err, avataroptiondomain.ErrAvatarOptionNotFound) {
		return huma.Error404NotFound("avatar option not found")
	}
	return err
}
