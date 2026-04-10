package authdomain

import "time"

type OAuthState struct {
	IssuedAt   time.Time
	Nonce      string
	RedirectTo string
}

type OAuthStateCodec interface {
	Decode(rawState string) (*OAuthState, error)
	Encode(state OAuthState) (string, error)
}
