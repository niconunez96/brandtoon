package avatarusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatardto "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases/dto"
	"context"
)

type ListAvatarsQuery struct {
	UserID string
}

func ListAvatars(
	ctx context.Context,
	query ListAvatarsQuery,
	avatarRepo avatardomain.AvatarRepository,
) ([]avatardto.AvatarDTO, error) {
	avatars, err := avatarRepo.ListByUserID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	return serializeList(avatars), nil
}
