package mocks

import (
	"context"

	"brandtoonapi/bounded_contexts/identity/user/domain"
)

type UserRepositoryMock struct {
	CreateFunc      func(ctx context.Context, user domain.User) error
	FindByEmailFunc func(ctx context.Context, email string) (*domain.User, error)
	FindByIDFunc    func(ctx context.Context, id string) (*domain.User, error)
	UpdateFunc      func(ctx context.Context, user domain.User) error
}

func (m *UserRepositoryMock) Create(ctx context.Context, user domain.User) error {
	if m.CreateFunc == nil {
		return nil
	}

	return m.CreateFunc(ctx, user)
}

func (m *UserRepositoryMock) FindByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {
	if m.FindByEmailFunc == nil {
		return nil, nil
	}

	return m.FindByEmailFunc(ctx, email)
}

func (m *UserRepositoryMock) FindByID(ctx context.Context, id string) (*domain.User, error) {
	if m.FindByIDFunc == nil {
		return nil, nil
	}

	return m.FindByIDFunc(ctx, id)
}

func (m *UserRepositoryMock) Update(ctx context.Context, user domain.User) error {
	if m.UpdateFunc == nil {
		return nil
	}

	return m.UpdateFunc(ctx, user)
}
