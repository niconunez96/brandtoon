package avatarhttp

import (
	avatarusecases "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	"context"

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
			return nil, err
		}

		response := &listAvatarsOutput{}
		response.Body.Avatars = make([]avatarPayload, 0, len(avatars))
		for _, avatar := range avatars {
			response.Body.Avatars = append(response.Body.Avatars, toAvatarPayload(avatar))
		}

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
			return nil, err
		}

		response := &createAvatarOutput{}
		response.Body.Avatar = toAvatarPayload(avatar)
		return response, nil
	}
}

func requireAuthUserMetadata(ctx context.Context) (*shareddomain.AuthUserMetadata, error) {
	userMetadata, ok := ctx.Value(shareddomain.UserMetadataContextKey).(*shareddomain.AuthUserMetadata)
	if !ok || userMetadata == nil || userMetadata.UserId == "" {
		return nil, huma.Error401Unauthorized("missing or invalid session")
	}

	return userMetadata, nil
}
