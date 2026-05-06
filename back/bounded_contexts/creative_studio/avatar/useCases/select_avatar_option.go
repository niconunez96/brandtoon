package avatarusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
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
) (AvatarOptionsDTO, error) {
	avatar, err := avatarRepo.FindOwnedByID(ctx, cmd.AvatarID, cmd.UserID)
	if err != nil {
		return AvatarOptionsDTO{}, err
	}

	if avatar == nil {
		return AvatarOptionsDTO{}, ErrAvatarNotFound
	}

	selectedAvatar, err := avatar.SelectOption(cmd.ID)
	if err != nil {
		if errors.Is(err, avatardomain.ErrAvatarOptionNotFound) {
			return AvatarOptionsDTO{}, ErrAvatarOptionNotFound
		}

		return AvatarOptionsDTO{}, err
	}

	if err := avatarRepo.UpdateOptions(ctx, cmd.AvatarID, cmd.UserID, selectedAvatar.AvatarOptions); err != nil {
		return AvatarOptionsDTO{}, err
	}

	return serializeAvatarOptions(selectedAvatar), nil
}
