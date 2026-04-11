package userusecases

import (
	userdomain "brandtoonapi/bounded_contexts/identity/user/domain"
	"context"
)

type FindUserQuery struct {
	UserId string
	Email  string
}

func FindUser(
	ctx context.Context,
	query FindUserQuery,
	userRepo userdomain.UserRepository,
) (*userdomain.User, error) {
	if query.UserId != "" {
		return userRepo.FindByID(ctx, query.UserId)
	}
	return userRepo.FindByEmail(ctx, query.Email)
}
