package mocks

import (
	"context"

	avataroptiondomain "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain"
)

type AvatarOptionRepositoryMock struct {
	CreatePendingBatchFunc func(ctx context.Context, options []avataroptiondomain.AvatarOption) error
	ListByAvatarIDFunc     func(ctx context.Context, avatarID string) ([]avataroptiondomain.AvatarOption, error)
	SelectFunc             func(ctx context.Context, avatarID string, optionID string) error
	DeleteFunc             func(ctx context.Context, avatarID string, ids []string) error
	MarkDoneFunc           func(ctx context.Context, id string, href string) error
	MarkFailedFunc         func(ctx context.Context, id string) error
}

func (m *AvatarOptionRepositoryMock) CreatePendingBatch(
	ctx context.Context,
	options []avataroptiondomain.AvatarOption,
) error {
	if m.CreatePendingBatchFunc == nil {
		return nil
	}

	return m.CreatePendingBatchFunc(ctx, options)
}

func (m *AvatarOptionRepositoryMock) ListByAvatarID(
	ctx context.Context,
	avatarID string,
) ([]avataroptiondomain.AvatarOption, error) {
	if m.ListByAvatarIDFunc == nil {
		return nil, nil
	}

	return m.ListByAvatarIDFunc(ctx, avatarID)
}

func (m *AvatarOptionRepositoryMock) Select(ctx context.Context, avatarID string, optionID string) error {
	if m.SelectFunc == nil {
		return nil
	}

	return m.SelectFunc(ctx, avatarID, optionID)
}

func (m *AvatarOptionRepositoryMock) Delete(ctx context.Context, avatarID string, ids []string) error {
	if m.DeleteFunc == nil {
		return nil
	}

	return m.DeleteFunc(ctx, avatarID, ids)
}

func (m *AvatarOptionRepositoryMock) MarkDone(ctx context.Context, id string, href string) error {
	if m.MarkDoneFunc == nil {
		return nil
	}

	return m.MarkDoneFunc(ctx, id, href)
}

func (m *AvatarOptionRepositoryMock) MarkFailed(ctx context.Context, id string) error {
	if m.MarkFailedFunc == nil {
		return nil
	}

	return m.MarkFailedFunc(ctx, id)
}
