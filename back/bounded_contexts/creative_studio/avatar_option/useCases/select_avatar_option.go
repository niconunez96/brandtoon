package avataroptionusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avataroptiondomain "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain"
	avataroptiondto "brandtoonapi/bounded_contexts/creative_studio/avatar_option/useCases/dto"
	"context"
	"errors"
)

type SelectAvatarOptionCommand struct {
	AvatarID string
	ID       string
	UserID   string
}

func SelectAvatarOption(
	ctx context.Context,
	cmd SelectAvatarOptionCommand,
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

	if err := avatarOptionRepo.Select(ctx, cmd.AvatarID, cmd.ID); err != nil {
		if errors.Is(err, avataroptiondomain.ErrAvatarOptionNotFound) {
			return avataroptiondto.AvatarOptionsDTO{}, ErrAvatarOptionNotFound
		}

		return avataroptiondto.AvatarOptionsDTO{}, err
	}

	options, err := avatarOptionRepo.ListByAvatarID(ctx, cmd.AvatarID)
	if err != nil {
		return avataroptiondto.AvatarOptionsDTO{}, err
	}

	return avataroptiondto.SerializeList(cmd.AvatarID, options), nil
}
