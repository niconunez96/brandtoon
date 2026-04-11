package test

import (
	"context"
	"testing"

	sessionmocks "brandtoonapi/bounded_contexts/identity/session/domain/mocks"
	sessionusecases "brandtoonapi/bounded_contexts/identity/session/useCases"
)

func TestLogoutSessionDeletesSessionWhenPresent(t *testing.T) {
	t.Parallel()

	deletedSessionID := ""
	err := sessionusecases.DeleteSession(
		context.Background(),
		sessionusecases.LogoutSessionCommand{SessionID: "session-v7"},
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
