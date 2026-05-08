package avatarusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatardto "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases/dto"
	"context"
)

type GetAvatarQuery struct {
	AvatarID string
	UserID   string
}

func GetAvatar(
	ctx context.Context,
	query GetAvatarQuery,
	avatarRepo avatardomain.AvatarRepository,
) (avatardto.AvatarDetailsDTO, error) {
	avatar, err := avatarRepo.FindOwnedByID(ctx, query.AvatarID, query.UserID)
	if err != nil {
		return avatardto.AvatarDetailsDTO{}, err
	}

	if avatar == nil {
		return avatardto.AvatarDetailsDTO{}, ErrAvatarNotFound
	}

	return serializeAvatarDetails(avatar.WithNormalizedOptions()), nil
}
