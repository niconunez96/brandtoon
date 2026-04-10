package authusecases

import (
	"context"
	"fmt"
	"time"

	"brandtoonapi/bounded_contexts/identity/auth/domain"
	"brandtoonapi/bounded_contexts/identity/session/domain"
	"brandtoonapi/bounded_contexts/identity/user/domain"
)

type GetCurrentUserQuery struct {
	SessionID string
}

type GetCurrentUserResult struct {
	Session sessiondomain.Session
	User    userdomain.User
}

func GetCurrentUser(
	ctx context.Context,
	query GetCurrentUserQuery,
	sessionRepo sessiondomain.SessionRepository,
	userRepo userdomain.UserRepository,
	now func() time.Time,
) (GetCurrentUserResult, error) {
	if query.SessionID == "" {
		fmt.Println("SESSION ID IS EMPTY")
		return GetCurrentUserResult{}, authdomain.ErrUnauthenticated
	}

	session, err := sessionRepo.FindActiveByID(ctx, query.SessionID)
	if err != nil {
		return GetCurrentUserResult{}, err
	}

	if session == nil || session.IsExpired(now().UTC()) {
		fmt.Println("SESSION IS EXPIRED")
		return GetCurrentUserResult{}, authdomain.ErrUnauthenticated
	}

	user, err := userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return GetCurrentUserResult{}, err
	}

	if user == nil {
		fmt.Println("NO USER FOUND")
		return GetCurrentUserResult{}, authdomain.ErrUnauthenticated
	}

	return GetCurrentUserResult{Session: *session, User: *user}, nil
}
