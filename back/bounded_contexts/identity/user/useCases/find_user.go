package usecases

import (
	"brandtoonapi/bounded_contexts/identity/user/domain"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	"context"
)

type FindUserQuery struct {
	Email string
}

func FindUser(
	ctx context.Context,
	query FindUserQuery,
	userRepo domain.UserRepository,
	idGenerator shareddomain.IDGenerator,
) (*domain.User, error) {
	return userRepo.FindByEmail(ctx, query.Email)
}
