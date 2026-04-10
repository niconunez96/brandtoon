package shared

import (
	"context"
	"sync"

	"brandtoonapi/bounded_contexts/identity/auth/domain"
	authoauth "brandtoonapi/bounded_contexts/identity/auth/infra/oauth"
	authsecurity "brandtoonapi/bounded_contexts/identity/auth/infra/security"
	sessiondomain "brandtoonapi/bounded_contexts/identity/session/domain"
	sessionrepo "brandtoonapi/bounded_contexts/identity/session/infra/repo"
	userdomain "brandtoonapi/bounded_contexts/identity/user/domain"
	userrepo "brandtoonapi/bounded_contexts/identity/user/infra/repo"
	sharedconfig "brandtoonapi/bounded_contexts/shared/infra/config"
	sharedpostgres "brandtoonapi/bounded_contexts/shared/infra/postgres"

	"github.com/jmoiron/sqlx"
)

var (
	container     *DIContainer
	containerOnce sync.Once
)

type DIContainer struct {
	configOnce sync.Once
	config     sharedconfig.Config
	configErr  error

	dbOnce sync.Once
	db     *sqlx.DB
	dbErr  error

	googleProviderOnce sync.Once
	googleProvider     domain.IdentityProvider
	googleProviderErr  error

	stateCodecOnce sync.Once
	stateCodec     domain.OAuthStateCodec
	stateCodecErr  error

	userRepoOnce sync.Once
	userRepo     userdomain.UserRepository
	userRepoErr  error

	sessionRepoOnce sync.Once
	sessionRepo     sessiondomain.SessionRepository
	sessionRepoErr  error
}

func NewDIContainer() *DIContainer {
	containerOnce.Do(func() {
		container = &DIContainer{}
	})

	return container
}

func (c *DIContainer) GetConfig() (sharedconfig.Config, error) {
	c.configOnce.Do(func() {
		c.config, c.configErr = sharedconfig.LoadConfig()
	})

	return c.config, c.configErr
}

func (c *DIContainer) GetDB(ctx context.Context) (*sqlx.DB, error) {
	c.dbOnce.Do(func() {
		config, err := c.GetConfig()
		if err != nil {
			c.dbErr = err
			return
		}

		c.db, c.dbErr = sharedpostgres.NewDatabase(ctx, config.DatabaseURL)
	})

	return c.db, c.dbErr
}

func (c *DIContainer) GetGoogleIdentityProvider() (domain.IdentityProvider, error) {
	c.googleProviderOnce.Do(func() {
		config, err := c.GetConfig()
		if err != nil {
			c.googleProviderErr = err
			return
		}

		c.googleProvider = authoauth.NewGoogleOAuthClient(
			config.GoogleClientID,
			config.GoogleClientSecret,
			config.GoogleRedirectURL,
		)
	})

	return c.googleProvider, c.googleProviderErr
}

func (c *DIContainer) GetOAuthStateCodec() (domain.OAuthStateCodec, error) {
	c.stateCodecOnce.Do(func() {
		config, err := c.GetConfig()
		if err != nil {
			c.stateCodecErr = err
			return
		}

		c.stateCodec = authsecurity.NewHMACOAuthStateCodec(config.AuthStateSecret)
	})

	return c.stateCodec, c.stateCodecErr
}

func (c *DIContainer) GetUserRepo(ctx context.Context) (userdomain.UserRepository, error) {
	c.userRepoOnce.Do(func() {
		db, err := c.GetDB(ctx)
		if err != nil {
			c.userRepoErr = err
			return
		}

		c.userRepo = userrepo.NewUserPostgresRepo(db)
	})

	return c.userRepo, c.userRepoErr
}

func (c *DIContainer) GetSessionRepo(ctx context.Context) (sessiondomain.SessionRepository, error) {
	c.sessionRepoOnce.Do(func() {
		db, err := c.GetDB(ctx)
		if err != nil {
			c.sessionRepoErr = err
			return
		}

		c.sessionRepo = sessionrepo.NewSessionPostgresRepo(db)
	})

	return c.sessionRepo, c.sessionRepoErr
}
