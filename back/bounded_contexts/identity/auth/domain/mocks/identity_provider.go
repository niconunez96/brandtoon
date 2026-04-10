package mocks

import (
	authdomain "brandtoonapi/bounded_contexts/identity/auth/domain"
	"context"
)

type IdentityProviderMock struct {
	BuildAuthURLFunc func(state string) string
	ExchangeCodeFunc func(ctx context.Context, code string) (*authdomain.Identity, error)
}

func (m *IdentityProviderMock) BuildAuthURL(state string) string {
	if m.BuildAuthURLFunc == nil {
		return ""
	}

	return m.BuildAuthURLFunc(state)
}

func (m *IdentityProviderMock) ExchangeCode(
	ctx context.Context,
	code string,
) (*authdomain.Identity, error) {
	if m.ExchangeCodeFunc == nil {
		return nil, nil
	}

	return m.ExchangeCodeFunc(ctx, code)
}
