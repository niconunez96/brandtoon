package authusecases

import (
	authdomain "brandtoonapi/bounded_contexts/identity/auth/domain"
	sessiondomain "brandtoonapi/bounded_contexts/identity/session/domain"
	userdomain "brandtoonapi/bounded_contexts/identity/user/domain"
	userusecases "brandtoonapi/bounded_contexts/identity/user/useCases"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	"context"
	"time"
)

type AuthenticateCallbackCommand struct {
	Code       string
	SessionTTL time.Duration
	State      string
}

type AuthenticateCallbackResult struct {
	RedirectTo string
	Session    sessiondomain.Session
	User       userdomain.User
}

func AuthenticateCallback(
	ctx context.Context,
	command AuthenticateCallbackCommand,
	identityProvider authdomain.IdentityProvider,
	stateCodec authdomain.OAuthStateCodec,
	userRepo userdomain.UserRepository,
	sessionRepo sessiondomain.SessionRepository,
	idGenerator shareddomain.IDGenerator,
	now func() time.Time,
) (*AuthenticateCallbackResult, error) {
	state, err := stateCodec.Decode(command.State)
	if err != nil {
		return nil, err
	}

	identity, err := identityProvider.ExchangeCode(ctx, command.Code)
	if err != nil {
		return nil, err
	}

	user, err := findOrCreateUser(ctx, identity, userRepo, idGenerator)
	if err != nil {
		return nil, err
	}

	sessionID, err := idGenerator()
	if err != nil {
		return nil, err
	}

	session := sessiondomain.NewSession(sessionID, user.ID, now().UTC().Add(command.SessionTTL))
	if err := sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &AuthenticateCallbackResult{
		RedirectTo: sanitizeRedirectTarget(state.RedirectTo),
		Session:    session,
		User:       user,
	}, nil
}

func findOrCreateUser(
	ctx context.Context,
	identity *authdomain.Identity,
	userRepo userdomain.UserRepository,
	idGenerator shareddomain.IDGenerator,
) (userdomain.User, error) {
	existingUser, err := userusecases.FindUser(
		ctx,
		userusecases.FindUserQuery{Email: identity.Email},
		userRepo,
		idGenerator,
	)
	if err != nil {
		return userdomain.User{}, err
	}

	if existingUser == nil {
		return userusecases.CreateUser(ctx, userusecases.CreateUserCmd{
			Name:      identity.Name,
			Email:     identity.Email,
			AvatarURL: identity.AvatarURL,
		}, userRepo, idGenerator)
	}

	updatedUser := existingUser.UpdateProfile(identity.Email, identity.Name, identity.AvatarURL)
	if err := userRepo.Update(ctx, updatedUser); err != nil {
		return userdomain.User{}, err
	}

	return updatedUser, nil
}
