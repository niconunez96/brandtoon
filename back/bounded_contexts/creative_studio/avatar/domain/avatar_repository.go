package avatardomain

import "context"

type AvatarRepository interface {
	Create(ctx context.Context, avatar Avatar) error
	ListByUserID(ctx context.Context, userID string) ([]Avatar, error)
}
