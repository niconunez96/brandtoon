package avatarusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
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
) (AvatarOptionsDTO, error) {
	avatar, err := avatarRepo.FindOwnedByID(ctx, cmd.AvatarID, cmd.UserID)
	if err != nil {
		return AvatarOptionsDTO{}, err
	}

	if avatar == nil {
		return AvatarOptionsDTO{}, ErrAvatarNotFound
	}

	updatedAvatar, err := avatar.DeleteOptions(cmd.IDs)
	if err != nil {
		if errors.Is(err, avatardomain.ErrAvatarOptionNotFound) {
			return AvatarOptionsDTO{}, ErrAvatarOptionNotFound
		}

		return AvatarOptionsDTO{}, err
	}

	if err := avatarRepo.UpdateOptions(ctx, cmd.AvatarID, cmd.UserID, updatedAvatar.AvatarOptions); err != nil {
		return AvatarOptionsDTO{}, err
	}

	return serializeAvatarOptions(updatedAvatar), nil
}
