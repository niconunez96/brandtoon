package avatarhttp

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarusecases "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	"context"
	"errors"
	stdhttp "net/http"

	"github.com/danielgtaylor/huma/v2"
)

func buildListAvatarsHandler(
	deps RouteDependencies,
) func(ctx context.Context, input *struct{}) (*listAvatarsOutput, error) {
	return func(ctx context.Context, input *struct{}) (*listAvatarsOutput, error) {
		userMetadata, err := requireAuthUserMetadata(ctx)
		if err != nil {
			return nil, err
		}

		avatars, err := avatarusecases.ListAvatars(
			ctx,
			avatarusecases.ListAvatarsQuery{UserID: userMetadata.UserId},
			deps.AvatarRepo,
		)
		if err != nil {
			if errors.Is(err, avatardomain.ErrInvalidName) {
				return nil, huma.Error422UnprocessableEntity("invalid avatar name")
			}
			return nil, err
		}

		response := &listAvatarsOutput{}
		response.Body.Avatars = avatars

		return response, nil
	}
}

func buildCreateAvatarHandler(
	deps RouteDependencies,
) func(ctx context.Context, input *createAvatarInput) (*createAvatarOutput, error) {
	return func(ctx context.Context, input *createAvatarInput) (*createAvatarOutput, error) {
		userMetadata, err := requireAuthUserMetadata(ctx)
		if err != nil {
			return nil, err
		}

		avatar, err := avatarusecases.CreateAvatar(
			ctx,
			avatarusecases.CreateAvatarCommand{
				Name:   input.Body.Name,
				UserID: userMetadata.UserId,
			},
			deps.AvatarRepo,
			deps.IDGenerator,
		)
		if err != nil {
			if errors.Is(err, avatardomain.ErrInvalidName) {
				return nil, huma.Error422UnprocessableEntity("invalid avatar name")
			}
			return nil, err
		}

		response := &createAvatarOutput{}
		response.Body.Avatar = avatar
		return response, nil
	}
}

func buildGetAvatarHandler(
	deps RouteDependencies,
) func(ctx context.Context, input *getAvatarInput) (*getAvatarOutput, error) {
	return func(ctx context.Context, input *getAvatarInput) (*getAvatarOutput, error) {
		userMetadata, err := requireAuthUserMetadata(ctx)
		if err != nil {
			return nil, err
		}

		avatar, err := avatarusecases.GetAvatar(
			ctx,
			avatarusecases.GetAvatarQuery{
				AvatarID: input.AvatarID,
				UserID:   userMetadata.UserId,
			},
			deps.AvatarRepo,
		)
		if err != nil {
			if errors.Is(err, avatarusecases.ErrAvatarNotFound) {
				return nil, huma.Error404NotFound("avatar not found")
			}

			return nil, err
		}

		response := &getAvatarOutput{}
		response.Body.Avatar = avatar
		return response, nil
	}
}

func buildGenerateAvatarOptionsHandler(
	deps RouteDependencies,
) stdhttp.HandlerFunc {
	return func(writer stdhttp.ResponseWriter, request *stdhttp.Request) {
		userMetadata, err := requireAuthUserMetadata(request.Context())
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

func buildSelectAvatarOptionHandler(
	deps RouteDependencies,
) func(ctx context.Context, input *selectAvatarOptionInput) (*avatarOptionsOutput, error) {
	return func(ctx context.Context, input *selectAvatarOptionInput) (*avatarOptionsOutput, error) {
		userMetadata, err := requireAuthUserMetadata(ctx)
		if err != nil {
			return nil, err
		}

		avatar, err := avatarusecases.SelectAvatarOption(
			ctx,
			avatarusecases.SelectAvatarOptionCommand{
				AvatarID: input.AvatarID,
				ID:       input.Body.ID,
				UserID:   userMetadata.UserId,
			},
			deps.AvatarRepo,
		)
		if err != nil {
			return nil, mapAvatarOptionsError(err)
		}

		response := &avatarOptionsOutput{}
		response.Body.Avatar = avatar
		return response, nil
	}
}

func buildDeleteAvatarOptionsHandler(
	deps RouteDependencies,
) func(ctx context.Context, input *deleteAvatarOptionsInput) (*avatarOptionsOutput, error) {
	return func(ctx context.Context, input *deleteAvatarOptionsInput) (*avatarOptionsOutput, error) {
		userMetadata, err := requireAuthUserMetadata(ctx)
		if err != nil {
			return nil, err
		}

		avatar, err := avatarusecases.DeleteAvatarOptions(
			ctx,
			avatarusecases.DeleteAvatarOptionsCommand{
				AvatarID: input.AvatarID,
				IDs:      input.Body.IDs,
				UserID:   userMetadata.UserId,
			},
			deps.AvatarRepo,
		)
		if err != nil {
			return nil, mapAvatarOptionsError(err)
		}

		response := &avatarOptionsOutput{}
		response.Body.Avatar = avatar
		return response, nil
	}
}

func RequireAuthUserMetadata(ctx context.Context) (*shareddomain.AuthUserMetadata, error) {
	userMetadata, ok := ctx.Value(shareddomain.UserMetadataContextKey).(*shareddomain.AuthUserMetadata)
	if !ok || userMetadata == nil || userMetadata.UserId == "" {
		return nil, huma.Error401Unauthorized("missing or invalid session")
	}

	return userMetadata, nil
}

func requireAuthUserMetadata(ctx context.Context) (*shareddomain.AuthUserMetadata, error) {
	return RequireAuthUserMetadata(ctx)
}

func mapAvatarOptionsError(err error) error {
	if errors.Is(err, avatarusecases.ErrAvatarNotFound) {
		return huma.Error404NotFound("avatar not found")
	}

	if errors.Is(err, avatarusecases.ErrAvatarOptionNotFound) {
		return huma.Error404NotFound("avatar option not found")
	}

	return err
}
