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
) (AvatarDTO, error) {
	normalizedName := strings.TrimSpace(cmd.Name)
	if normalizedName == "" {
		return AvatarDTO{}, avatardomain.ErrInvalidName
	}

	avatarID, err := idGenerator()
	if err != nil {
		return AvatarDTO{}, err
	}

	avatar := avatardomain.NewAvatar(avatarID, cmd.UserID, normalizedName)
	if err := avatarRepo.Create(ctx, avatar); err != nil {
		return AvatarDTO{}, err
	}

	return serialize(avatar), nil
}
