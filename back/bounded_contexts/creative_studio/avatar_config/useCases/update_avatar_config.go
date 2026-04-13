package avatarconfigusecases

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	"context"
)

type UpdateAvatarConfigCommand struct {
	AvatarID      string
	ArtisticStyle string
	Prompt        string
	UserID        string
}

func UpdateAvatarConfig(
	ctx context.Context,
	cmd UpdateAvatarConfigCommand,
	avatarRepo avatardomain.AvatarRepository,
	avatarConfigRepo avatarconfigdomain.AvatarConfigRepository,
) (avatarconfigdomain.AvatarConfig, error) {
	avatar, err := avatarRepo.FindOwnedByID(ctx, cmd.AvatarID, cmd.UserID)
	if err != nil {
		return avatarconfigdomain.AvatarConfig{}, err
	}

	if avatar == nil {
		return avatarconfigdomain.AvatarConfig{}, ErrAvatarNotFound
	}

	artisticStyle, err := avatarconfigdomain.ParseArtisticStyle(cmd.ArtisticStyle)
	if err != nil {
		return avatarconfigdomain.AvatarConfig{}, err
	}

	avatarConfig := avatarconfigdomain.NewAvatarConfig(
		cmd.AvatarID,
		cmd.Prompt,
		artisticStyle,
	)
	if err := avatarConfigRepo.Upsert(ctx, avatarConfig); err != nil {
		return avatarconfigdomain.AvatarConfig{}, err
	}

	return avatarConfig, nil
}
