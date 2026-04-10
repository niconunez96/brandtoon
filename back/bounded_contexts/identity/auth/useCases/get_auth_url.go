package authusecases

import (
	authdomain "brandtoonapi/bounded_contexts/identity/auth/domain"
	"time"
)

type GetAuthURLQuery struct {
	RedirectTo string
}

type AuthURLResult struct {
	URL string
}

func GetAuthURL(
	query GetAuthURLQuery,
	stateCodec authdomain.OAuthStateCodec,
	authProvider authdomain.IdentityProvider,
	now time.Time,
) (*AuthURLResult, error) {
	state, err := stateCodec.Encode(authdomain.OAuthState{
		IssuedAt:   now.UTC(),
		Nonce:      now.UTC().Format(time.RFC3339Nano),
		RedirectTo: sanitizeRedirectTarget(query.RedirectTo),
	})
	if err != nil {
		return nil, err
	}

	return &AuthURLResult{URL: authProvider.BuildAuthURL(state)}, nil
}
