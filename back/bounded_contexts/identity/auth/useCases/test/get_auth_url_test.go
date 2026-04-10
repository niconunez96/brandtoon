package test

import (
	"testing"
	"time"

	authdomain "brandtoonapi/bounded_contexts/identity/auth/domain"
	authmocks "brandtoonapi/bounded_contexts/identity/auth/domain/mocks"
	authusecases "brandtoonapi/bounded_contexts/identity/auth/useCases"
)

func TestGetGoogleAuthURLDefaultsInvalidRedirectTarget(t *testing.T) {
	t.Parallel()

	var encodedState authdomain.OAuthState
	stateCodec := &authmocks.OAuthStateCodecMock{
		EncodeFunc: func(state authdomain.OAuthState) (string, error) {
			encodedState = state
			return "signed-state", nil
		},
	}
	provider := &authmocks.IdentityProviderMock{
		BuildAuthURLFunc: func(state string) string {
			if state != "signed-state" {
				t.Fatalf("expected signed-state, got %s", state)
			}

			return "https://accounts.google.com/o/oauth2/auth"
		},
	}

	result, err := authusecases.GetAuthURL(
		authusecases.GetAuthURLQuery{RedirectTo: "https://malicious.example.com"},
		stateCodec,
		provider,
		time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if result.URL == "" {
		t.Fatal("expected auth url to be returned")
	}

	if encodedState.RedirectTo != "/creative-studio" {
		t.Fatalf("expected default redirect target, got %s", encodedState.RedirectTo)
	}
}
