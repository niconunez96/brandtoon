package test

import (
	"context"
	"testing"

	"brandtoonapi/bounded_contexts/identity/auth/useCases"
	sessionmocks "brandtoonapi/bounded_contexts/identity/session/domain/mocks"
)

func TestLogoutSessionDeletesSessionWhenPresent(t *testing.T) {
	t.Parallel()

	deletedSessionID := ""
	err := authusecases.LogoutSession(
		context.Background(),
		authusecases.LogoutSessionCommand{SessionID: "session-v7"},
		&sessionmocks.SessionRepositoryMock{
			DeleteFunc: func(ctx context.Context, id string) error {
				deletedSessionID = id
				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if deletedSessionID != "session-v7" {
		t.Fatalf("expected session-v7, got %s", deletedSessionID)
	}
}
