package avatarhttp

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarusecases "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	"context"
	"errors"

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
