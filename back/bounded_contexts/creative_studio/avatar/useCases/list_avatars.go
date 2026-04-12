package avatarusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	"context"
)

type ListAvatarsQuery struct {
	UserID string
}

func ListAvatars(
	ctx context.Context,
	query ListAvatarsQuery,
	avatarRepo avatardomain.AvatarRepository,
) ([]avatardomain.Avatar, error) {
	return avatarRepo.ListByUserID(ctx, query.UserID)
}
