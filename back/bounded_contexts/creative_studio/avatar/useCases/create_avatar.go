package avatarusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	"context"
	"strings"
)

type CreateAvatarCommand struct {
	Name   string
	UserID string
}

func CreateAvatar(
	ctx context.Context,
	cmd CreateAvatarCommand,
	avatarRepo avatardomain.AvatarRepository,
	idGenerator shareddomain.IDGenerator,
) (avatardomain.Avatar, error) {
	normalizedName := strings.TrimSpace(cmd.Name)

	avatarID, err := idGenerator()
	if err != nil {
		return avatardomain.Avatar{}, err
	}

	avatar := avatardomain.NewAvatar(avatarID, cmd.UserID, normalizedName)
	if err := avatarRepo.Create(ctx, avatar); err != nil {
		return avatardomain.Avatar{}, err
	}

	return avatar, nil
}
