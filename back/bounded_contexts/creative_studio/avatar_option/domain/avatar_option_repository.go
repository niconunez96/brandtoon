package avataroptiondomain

import "context"

type AvatarOptionRepository interface {
	CreatePendingBatch(ctx context.Context, options []AvatarOption) error
	ListByAvatarID(ctx context.Context, avatarID string) ([]AvatarOption, error)
	Select(ctx context.Context, avatarID string, optionID string) error
	Delete(ctx context.Context, avatarID string, ids []string) error
	MarkDone(ctx context.Context, id string, href string) error
	MarkFailed(ctx context.Context, id string) error
}
