package usecases

import (
	"brandtoonapi/bounded_contexts/identity/user/domain"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
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
	userRepo domain.UserRepository,
	idGenerator shareddomain.IDGenerator,
) (domain.User, error) {
	existingUser, err := userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return domain.User{}, err
	}
	if existingUser != nil {
		return domain.User{}, errors.New("user already exist")
	}
	userID, err := idGenerator()
	if err != nil {
		return domain.User{}, err
	}

	user := domain.NewUser(
		userID,
		cmd.Email,
		cmd.Name,
		cmd.AvatarURL,
	)
	if err := userRepo.Create(ctx, user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}
