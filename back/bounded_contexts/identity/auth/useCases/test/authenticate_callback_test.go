package test

import (
	"context"
	"testing"
	"time"

	authdomain "brandtoonapi/bounded_contexts/identity/auth/domain"
	authmocks "brandtoonapi/bounded_contexts/identity/auth/domain/mocks"
	usecases "brandtoonapi/bounded_contexts/identity/auth/useCases"
	sessiondomain "brandtoonapi/bounded_contexts/identity/session/domain"
	sessionmocks "brandtoonapi/bounded_contexts/identity/session/domain/mocks"
	userdomain "brandtoonapi/bounded_contexts/identity/user/domain"
	usermocks "brandtoonapi/bounded_contexts/identity/user/domain/mocks"
)

func TestAuthenticateGoogleCallbackCreatesUserAndSession(t *testing.T) {
	t.Parallel()

	createdUser := userdomain.User{}
	createdSession := sessiondomain.Session{}
	idCalls := 0
	now := time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC)

	result, err := usecases.AuthenticateCallback(
		context.Background(),
		usecases.AuthenticateCallbackCommand{
			Code:       "auth-code",
			SessionTTL: 30 * 24 * time.Hour,
			State:      "signed-state",
		},
		&authmocks.IdentityProviderMock{
			ExchangeCodeFunc: func(ctx context.Context, code string) (*authdomain.Identity, error) {
				if code != "auth-code" {
					t.Fatalf("expected auth-code, got %s", code)
				}

				return &authdomain.Identity{
					AvatarURL: "https://avatar.example.com/me.png",
					Email:     "nico@example.com",
					Name:      "Nico",
					Subject:   "google-123",
				}, nil
			},
		},
		&authmocks.OAuthStateCodecMock{
			DecodeFunc: func(rawState string) (*authdomain.OAuthState, error) {
				if rawState != "signed-state" {
					t.Fatalf("expected signed-state, got %s", rawState)
				}

				return &authdomain.OAuthState{RedirectTo: "/creative-studio"}, nil
			},
		},
		&usermocks.UserRepositoryMock{
			CreateFunc: func(ctx context.Context, user userdomain.User) error {
				createdUser = user
				return nil
			},
			FindByEmailFunc: func(ctx context.Context, email string) (*userdomain.User, error) {
				if email != "nico@example.com" {
					t.Fatalf("expected google-123, got %s", email)
				}

				return nil, nil
			},
		},
		&sessionmocks.SessionRepositoryMock{
			CreateFunc: func(ctx context.Context, session sessiondomain.Session) error {
				createdSession = session
				return nil
			},
		},
		func() (string, error) {
			idCalls++
			if idCalls == 1 {
				return "user-v7", nil
			}

			return "session-v7", nil
		},
		func() time.Time { return now },
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if result.RedirectTo != "/creative-studio" {
		t.Fatalf("expected /creative-studio redirect, got %s", result.RedirectTo)
	}

	if createdUser.ID != "user-v7" {
		t.Fatalf("expected user-v7 id, got %s", createdUser.ID)
	}

	if createdSession.ID != "session-v7" {
		t.Fatalf("expected session-v7 id, got %s", createdSession.ID)
	}

	if createdSession.ExpiresAt != now.Add(30*24*time.Hour) {
		t.Fatalf("unexpected session expiry: %s", createdSession.ExpiresAt)
	}
}
