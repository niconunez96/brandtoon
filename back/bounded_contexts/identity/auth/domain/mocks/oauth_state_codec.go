package mocks

import "brandtoonapi/bounded_contexts/identity/auth/domain"

type OAuthStateCodecMock struct {
	DecodeFunc func(rawState string) (*domain.OAuthState, error)
	EncodeFunc func(state domain.OAuthState) (string, error)
}

func (m *OAuthStateCodecMock) Decode(rawState string) (*domain.OAuthState, error) {
	if m.DecodeFunc == nil {
		return nil, nil
	}

	return m.DecodeFunc(rawState)
}

func (m *OAuthStateCodecMock) Encode(state domain.OAuthState) (string, error) {
	if m.EncodeFunc == nil {
		return "", nil
	}

	return m.EncodeFunc(state)
}
