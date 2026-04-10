package authdomain

import "context"

type Identity struct {
	AvatarURL string
	Email     string
	Name      string
	Subject   string
}

type IdentityProvider interface {
	BuildAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*Identity, error)
}
