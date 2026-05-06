package avatarusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
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
) (AvatarDetailsDTO, error) {
	avatar, err := avatarRepo.FindOwnedByID(ctx, query.AvatarID, query.UserID)
	if err != nil {
		return AvatarDetailsDTO{}, err
	}

	if avatar == nil {
		return AvatarDetailsDTO{}, ErrAvatarNotFound
	}

	return serializeAvatarDetails(avatar.WithNormalizedOptions()), nil
}
