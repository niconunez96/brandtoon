package shared

import (
	"context"
	"sync"

	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarrepo "brandtoonapi/bounded_contexts/creative_studio/avatar/infra/repo"
	authdomain "brandtoonapi/bounded_contexts/identity/auth/domain"
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

	// Auth
	googleProvider authdomain.IdentityProvider
	stateCodec     authdomain.OAuthStateCodec
	userRepo       userdomain.UserRepository
	sessionRepo    sessiondomain.SessionRepository
	// Creative studio
	avatarRepo avatardomain.AvatarRepository
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

func (c *DIContainer) GetGoogleIdentityProvider() (authdomain.IdentityProvider, error) {
	if c.googleProvider == nil {
		config, err := c.GetConfig()
		if err != nil {
			return nil, err
		}
		c.googleProvider = authoauth.NewGoogleOAuthClient(
			config.GoogleClientID,
			config.GoogleClientSecret,
			config.GoogleRedirectURL,
		)
	}
	return c.googleProvider, nil
}

func (c *DIContainer) GetOAuthStateCodec() (authdomain.OAuthStateCodec, error) {
	if c.stateCodec == nil {
		config, err := c.GetConfig()
		if err != nil {
			return nil, err
		}
		c.stateCodec = authsecurity.NewHMACOAuthStateCodec(config.AuthStateSecret)
	}
	return c.stateCodec, nil
}

func (c *DIContainer) GetUserRepo(ctx context.Context) (userdomain.UserRepository, error) {
	if c.userRepo == nil {
		db, err := c.GetDB(ctx)
		if err != nil {
			return nil, err
		}

		c.userRepo = userrepo.NewUserPostgresRepo(db)
	}
	return c.userRepo, nil
}

func (c *DIContainer) GetSessionRepo(ctx context.Context) (sessiondomain.SessionRepository, error) {
	if c.sessionRepo == nil {
		db, err := c.GetDB(ctx)
		if err != nil {
			return nil, err
		}

		c.sessionRepo = sessionrepo.NewSessionPostgresRepo(db)
	}

	return c.sessionRepo, nil
}

func (c *DIContainer) GetAvatarRepo(ctx context.Context) (avatardomain.AvatarRepository, error) {
	if c.avatarRepo == nil {
		db, err := c.GetDB(ctx)
		if err != nil {
			return nil, err
		}

		c.avatarRepo = avatarrepo.NewAvatarPostgresRepo(db)
	}

	return c.avatarRepo, nil
}
