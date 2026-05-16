package avataroptionusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avataroptiondomain "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain"
	avataroptiondto "brandtoonapi/bounded_contexts/creative_studio/avatar_option/useCases/dto"
	"context"
)

type ListAvatarOptionsQuery struct {
	AvatarID string
	UserID   string
}

func ListAvatarOptions(
	ctx context.Context,
	query ListAvatarOptionsQuery,
	avatarRepo avatardomain.AvatarRepository,
	avatarOptionRepo avataroptiondomain.AvatarOptionRepository,
) (avataroptiondto.AvatarOptionsDTO, error) {
	avatar, err := avatarRepo.FindOwnedByID(ctx, query.AvatarID, query.UserID)
	if err != nil {
		return avataroptiondto.AvatarOptionsDTO{}, err
	}

	if avatar == nil {
		return avataroptiondto.AvatarOptionsDTO{}, ErrAvatarNotFound
	}

	options, err := avatarOptionRepo.ListByAvatarID(ctx, query.AvatarID)
	if err != nil {
		return avataroptiondto.AvatarOptionsDTO{}, err
	}

	return avataroptiondto.SerializeList(query.AvatarID, options), nil
}
