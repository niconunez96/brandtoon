package mocks

import authdomain "brandtoonapi/bounded_contexts/identity/auth/domain"

type OAuthStateCodecMock struct {
	DecodeFunc func(rawState string) (*authdomain.OAuthState, error)
	EncodeFunc func(state authdomain.OAuthState) (string, error)
}

func (m *OAuthStateCodecMock) Decode(rawState string) (*authdomain.OAuthState, error) {
	if m.DecodeFunc == nil {
		return nil, nil
	}

	return m.DecodeFunc(rawState)
}

func (m *OAuthStateCodecMock) Encode(state authdomain.OAuthState) (string, error) {
	if m.EncodeFunc == nil {
		return "", nil
	}

	return m.EncodeFunc(state)
}
