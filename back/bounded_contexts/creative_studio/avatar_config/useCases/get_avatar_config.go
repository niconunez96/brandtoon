package avatarconfigusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	"context"
)

type GetAvatarConfigQuery struct {
	AvatarID string
	UserID   string
}

func GetAvatarConfig(
	ctx context.Context,
	query GetAvatarConfigQuery,
	avatarRepo avatardomain.AvatarRepository,
	avatarConfigRepo avatarconfigdomain.AvatarConfigRepository,
) (*AvatarConfigDTO, error) {
	avatar, err := avatarRepo.FindOwnedByID(ctx, query.AvatarID, query.UserID)
	if err != nil {
		return nil, err
	}

	if avatar == nil {
		return nil, ErrAvatarNotFound
	}

	avatarConfig, err := avatarConfigRepo.FindByAvatarID(ctx, query.AvatarID)
	if err != nil {
		return nil, err
	}

	if avatarConfig == nil {
		return nil, nil
	}

	normalizedAvatar := avatar.WithNormalizedOptions()
	serialized := serialize(*avatarConfig, normalizedAvatar)
	return &serialized, nil
}
