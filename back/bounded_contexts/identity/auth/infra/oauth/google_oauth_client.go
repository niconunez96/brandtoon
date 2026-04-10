package oauth

import (
	authdomain "brandtoonapi/bounded_contexts/identity/auth/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

type GoogleOAuthClient struct {
	config *oauth2.Config
}

type googleUserInfoResponse struct {
	Email   string `json:"email"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func NewGoogleOAuthClient(
	clientID string,
	clientSecret string,
	redirectURL string,
) *GoogleOAuthClient {
	return &GoogleOAuthClient{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"openid",
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
		},
	}
}

func (c *GoogleOAuthClient) BuildAuthURL(state string) string {
	return c.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (c *GoogleOAuthClient) ExchangeCode(
	ctx context.Context,
	code string,
) (*authdomain.Identity, error) {
	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := c.config.Client(ctx, token)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, googleUserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo returned status %d", response.StatusCode)
	}

	var payload googleUserInfoResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &authdomain.Identity{
		AvatarURL: payload.Picture,
		Email:     payload.Email,
		Name:      payload.Name,
		Subject:   payload.ID,
	}, nil
}
