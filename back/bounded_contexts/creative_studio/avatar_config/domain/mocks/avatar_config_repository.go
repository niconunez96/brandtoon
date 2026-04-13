package mocks

import (
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	"context"
)

type AvatarConfigRepositoryMock struct {
	FindByAvatarIDFunc func(ctx context.Context, avatarID string) (*avatarconfigdomain.AvatarConfig, error)
	UpsertFunc         func(ctx context.Context, avatarConfig avatarconfigdomain.AvatarConfig) error
}

func (m *AvatarConfigRepositoryMock) FindByAvatarID(
	ctx context.Context,
	avatarID string,
) (*avatarconfigdomain.AvatarConfig, error) {
	if m.FindByAvatarIDFunc == nil {
		return nil, nil
	}

	return m.FindByAvatarIDFunc(ctx, avatarID)
}

func (m *AvatarConfigRepositoryMock) Upsert(
	ctx context.Context,
	avatarConfig avatarconfigdomain.AvatarConfig,
) error {
	if m.UpsertFunc == nil {
		return nil
	}

	return m.UpsertFunc(ctx, avatarConfig)
}
