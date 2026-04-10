package userusecases

import (
	"brandtoonapi/bounded_contexts/identity/user/domain"
	"brandtoonapi/bounded_contexts/shared/domain"
	"context"
	"errors"
)

type CreateUserCmd struct {
	Name      string
	Email     string
	AvatarURL string
}

func CreateUser(
	ctx context.Context,
	cmd CreateUserCmd,
	userRepo userdomain.UserRepository,
	idGenerator shareddomain.IDGenerator,
) (userdomain.User, error) {
	existingUser, err := userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return userdomain.User{}, err
	}
	if existingUser != nil {
		return *existingUser, nil
	}
	userID, err := idGenerator()
	if err != nil {
		return userdomain.User{}, err
	}

	user := userdomain.NewUser(
		userID,
		cmd.Email,
		cmd.Name,
		cmd.AvatarURL,
	)
	if err := userRepo.Create(ctx, user); err != nil {
		return userdomain.User{}, err
	}

	return user, nil
}
