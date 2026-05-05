package mocks

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	"context"
)

type AvatarRepositoryMock struct {
	CreateFunc        func(ctx context.Context, avatar avatardomain.Avatar) error
	FindOwnedByIDFunc func(ctx context.Context, avatarID string, userID string) (*avatardomain.Avatar, error)
	ListByUserIDFunc  func(ctx context.Context, userID string) ([]avatardomain.Avatar, error)
	UpdateOptionsFunc func(
		ctx context.Context,
		avatarID string,
		userID string,
		options []avatardomain.AvatarOption,
	) error
}

func (m *AvatarRepositoryMock) Create(ctx context.Context, avatar avatardomain.Avatar) error {
	if m.CreateFunc == nil {
		return nil
	}

	return m.CreateFunc(ctx, avatar)
}

func (m *AvatarRepositoryMock) ListByUserID(
	ctx context.Context,
	userID string,
) ([]avatardomain.Avatar, error) {
	if m.ListByUserIDFunc == nil {
		return nil, nil
	}

	return m.ListByUserIDFunc(ctx, userID)
}

func (m *AvatarRepositoryMock) FindOwnedByID(
	ctx context.Context,
	avatarID string,
	userID string,
) (*avatardomain.Avatar, error) {
	if m.FindOwnedByIDFunc == nil {
		return nil, nil
	}

	return m.FindOwnedByIDFunc(ctx, avatarID, userID)
}

func (m *AvatarRepositoryMock) UpdateOptions(
	ctx context.Context,
	avatarID string,
	userID string,
	options []avatardomain.AvatarOption,
) error {
	if m.UpdateOptionsFunc == nil {
		return nil
	}

	return m.UpdateOptionsFunc(ctx, avatarID, userID, options)
}
