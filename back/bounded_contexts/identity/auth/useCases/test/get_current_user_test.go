package test

import (
	"context"
	"errors"
	"testing"
	"time"

	authdomain "brandtoonapi/bounded_contexts/identity/auth/domain"
	usecases "brandtoonapi/bounded_contexts/identity/auth/useCases"
	sessiondomain "brandtoonapi/bounded_contexts/identity/session/domain"
	sessionmocks "brandtoonapi/bounded_contexts/identity/session/domain/mocks"
	userdomain "brandtoonapi/bounded_contexts/identity/user/domain"
	usermocks "brandtoonapi/bounded_contexts/identity/user/domain/mocks"
)

func TestGetCurrentUserReturnsUnauthenticatedWithoutSession(t *testing.T) {
	t.Parallel()

	_, err := usecases.GetCurrentUser(
		context.Background(),
		usecases.GetCurrentUserQuery{},
		&sessionmocks.SessionRepositoryMock{},
		&usermocks.UserRepositoryMock{},
		time.Now,
	)
	if !errors.Is(err, authdomain.ErrUnauthenticated) {
		t.Fatalf("expected unauthenticated error, got %v", err)
	}
}

func TestGetCurrentUserReturnsUserWhenSessionExists(t *testing.T) {
	t.Parallel()

	result, err := usecases.GetCurrentUser(
		context.Background(),
		usecases.GetCurrentUserQuery{SessionID: "session-v7"},
		&sessionmocks.SessionRepositoryMock{
			FindActiveByIDFunc: func(ctx context.Context, id string) (*sessiondomain.Session, error) {
				return &sessiondomain.Session{
					ExpiresAt: time.Date(2026, time.April, 20, 12, 0, 0, 0, time.UTC),
					ID:        id,
					UserID:    "user-v7",
				}, nil
			},
		},
		&usermocks.UserRepositoryMock{
			FindByIDFunc: func(ctx context.Context, id string) (*userdomain.User, error) {
				user := userdomain.NewUser(id, "nico@example.com", "Nico", "https://avatar")
				return &user, nil
			},
		},
		func() time.Time {
			return time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC)
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if result.User.ID != "user-v7" {
		t.Fatalf("expected user-v7, got %s", result.User.ID)
	}
}
