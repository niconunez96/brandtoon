package authusecases

import (
	sessiondomain "brandtoonapi/bounded_contexts/identity/session/domain"
	"context"
)

type LogoutSessionCommand struct {
	SessionID string
}

func LogoutSession(
	ctx context.Context,
	command LogoutSessionCommand,
	sessionRepo sessiondomain.SessionRepository,
) error {
	if command.SessionID == "" {
		return nil
	}

	return sessionRepo.Delete(ctx, command.SessionID)
}
