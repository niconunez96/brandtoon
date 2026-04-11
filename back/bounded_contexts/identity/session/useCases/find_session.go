package sessionusecases

import (
	sessiondomain "brandtoonapi/bounded_contexts/identity/session/domain"
	"context"
)

type FindSessionQuery struct {
	SessionId string
}

func FindSession(
	ctx context.Context,
	query FindSessionQuery,
	sessionRepo sessiondomain.SessionRepository,
) (*sessiondomain.Session, error) {
	return sessionRepo.FindActiveByID(ctx, query.SessionId)
}
