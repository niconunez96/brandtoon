package avataroptionusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avataroptiondomain "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain"
	avataroptiondto "brandtoonapi/bounded_contexts/creative_studio/avatar_option/useCases/dto"
	"context"
)

type DeleteAvatarOptionsCommand struct {
	AvatarID string
	IDs      []string
	UserID   string
}

func DeleteAvatarOptions(
	ctx context.Context,
	cmd DeleteAvatarOptionsCommand,
	avatarRepo avatardomain.AvatarRepository,
	avatarOptionRepo avataroptiondomain.AvatarOptionRepository,
) (avataroptiondto.AvatarOptionsDTO, error) {
	avatar, err := avatarRepo.FindOwnedByID(ctx, cmd.AvatarID, cmd.UserID)
	if err != nil {
		return avataroptiondto.AvatarOptionsDTO{}, err
	}

	if avatar == nil {
		return avataroptiondto.AvatarOptionsDTO{}, ErrAvatarNotFound
	}

	if err := avatarOptionRepo.Delete(ctx, cmd.AvatarID, cmd.IDs); err != nil {
		return avataroptiondto.AvatarOptionsDTO{}, err
	}

	options, err := avatarOptionRepo.ListByAvatarID(ctx, cmd.AvatarID)
	if err != nil {
		return avataroptiondto.AvatarOptionsDTO{}, err
	}

	return avataroptiondto.SerializeList(cmd.AvatarID, options), nil
}
