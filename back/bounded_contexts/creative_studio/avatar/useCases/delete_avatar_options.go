package avatarusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatardto "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases/dto"
	"context"
	"errors"
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
) (avatardto.AvatarOptionsDTO, error) {
	avatar, err := avatarRepo.FindOwnedByID(ctx, cmd.AvatarID, cmd.UserID)
	if err != nil {
		return avatardto.AvatarOptionsDTO{}, err
	}

	if avatar == nil {
		return avatardto.AvatarOptionsDTO{}, ErrAvatarNotFound
	}

	updatedAvatar, err := avatar.DeleteOptions(cmd.IDs)
	if err != nil {
		if errors.Is(err, avatardomain.ErrAvatarOptionNotFound) {
			return avatardto.AvatarOptionsDTO{}, ErrAvatarOptionNotFound
		}

		return avatardto.AvatarOptionsDTO{}, err
	}

	if err := avatarRepo.UpdateOptions(ctx, cmd.AvatarID, cmd.UserID, updatedAvatar.AvatarOptions); err != nil {
		return avatardto.AvatarOptionsDTO{}, err
	}

	return serializeAvatarOptions(updatedAvatar), nil
}
