package avatarconfigdomain

import "context"

type AvatarConfigRepository interface {
	FindByAvatarID(ctx context.Context, avatarID string) (*AvatarConfig, error)
	Upsert(ctx context.Context, avatarConfig AvatarConfig) error
}
