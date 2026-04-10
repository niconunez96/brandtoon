package mocks

import (
	sessiondomain "brandtoonapi/bounded_contexts/identity/session/domain"
	"context"
)

type SessionRepositoryMock struct {
	CreateFunc         func(ctx context.Context, session sessiondomain.Session) error
	DeleteFunc         func(ctx context.Context, id string) error
	FindActiveByIDFunc func(ctx context.Context, id string) (*sessiondomain.Session, error)
}

func (m *SessionRepositoryMock) Create(ctx context.Context, session sessiondomain.Session) error {
	if m.CreateFunc == nil {
		return nil
	}

	return m.CreateFunc(ctx, session)
}

func (m *SessionRepositoryMock) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc == nil {
		return nil
	}

	return m.DeleteFunc(ctx, id)
}

func (m *SessionRepositoryMock) FindActiveByID(
	ctx context.Context,
	id string,
) (*sessiondomain.Session, error) {
	if m.FindActiveByIDFunc == nil {
		return nil, nil
	}

	return m.FindActiveByIDFunc(ctx, id)
}
