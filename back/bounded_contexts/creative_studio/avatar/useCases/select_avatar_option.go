package avatarusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatardto "brandtoonapi/bounded_contexts/creative_studio/avatar/useCases/dto"
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
) (avatardto.AvatarOptionsDTO, error) {
	avatar, err := avatarRepo.FindOwnedByID(ctx, cmd.AvatarID, cmd.UserID)
	if err != nil {
		return avatardto.AvatarOptionsDTO{}, err
	}

	if avatar == nil {
		return avatardto.AvatarOptionsDTO{}, ErrAvatarNotFound
	}

	selectedAvatar, err := avatar.SelectOption(cmd.ID)
	if err != nil {
		if errors.Is(err, avatardomain.ErrAvatarOptionNotFound) {
			return avatardto.AvatarOptionsDTO{}, ErrAvatarOptionNotFound
		}

		return avatardto.AvatarOptionsDTO{}, err
	}

	if err := avatarRepo.UpdateOptions(ctx, cmd.AvatarID, cmd.UserID, selectedAvatar.AvatarOptions); err != nil {
		return avatardto.AvatarOptionsDTO{}, err
	}

	return serializeAvatarOptions(selectedAvatar), nil
}
