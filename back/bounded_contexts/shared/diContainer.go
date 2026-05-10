package shared

import (
	"context"
	"sync"

	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	avatarrepo "brandtoonapi/bounded_contexts/creative_studio/avatar/infra/repo"
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	avatarconfigrepo "brandtoonapi/bounded_contexts/creative_studio/avatar_config/infra/repo"
	avataroptiondomain "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain"
	avataroptiongenerator "brandtoonapi/bounded_contexts/creative_studio/avatar_option/infra/generator"
	avataroptionrepo "brandtoonapi/bounded_contexts/creative_studio/avatar_option/infra/repo"
	avataroptionusecases "brandtoonapi/bounded_contexts/creative_studio/avatar_option/useCases"
	authdomain "brandtoonapi/bounded_contexts/identity/auth/domain"
	authoauth "brandtoonapi/bounded_contexts/identity/auth/infra/oauth"
	authsecurity "brandtoonapi/bounded_contexts/identity/auth/infra/security"
	sessiondomain "brandtoonapi/bounded_contexts/identity/session/domain"
	sessionrepo "brandtoonapi/bounded_contexts/identity/session/infra/repo"
	userdomain "brandtoonapi/bounded_contexts/identity/user/domain"
	userrepo "brandtoonapi/bounded_contexts/identity/user/infra/repo"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	sharedconfig "brandtoonapi/bounded_contexts/shared/infra/config"
	sharedevents "brandtoonapi/bounded_contexts/shared/infra/events"
	sharedpostgres "brandtoonapi/bounded_contexts/shared/infra/postgres"
	sharedsse "brandtoonapi/bounded_contexts/shared/infra/sse"
	sharedstorage "brandtoonapi/bounded_contexts/shared/infra/storage"

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
	avatarConfigRepo avatarconfigdomain.AvatarConfigRepository
	avatarGenerator  avataroptiondomain.AvatarGenerator
	avatarRepo       avatardomain.AvatarRepository
	avatarOptionRepo avataroptiondomain.AvatarOptionRepository
	eventBus         shareddomain.EventBus
	fileStorage      shareddomain.FileStorage
	sseConnector     *sharedsse.Connector
	sseHub           *sharedsse.Hub
}

func NewDIContainer() *DIContainer {
	containerOnce.Do(func() {
		container = &DIContainer{}
	})

	return container
}

func (c *DIContainer) GetEventBus() shareddomain.EventBus {
	if c.eventBus == nil {
		c.eventBus = sharedevents.NewInMemoryEventBus()
	}

	return c.eventBus
}

func (c *DIContainer) GetSSEHub() *sharedsse.Hub {
	if c.sseHub == nil {
		c.sseHub = sharedsse.NewHub()
	}

	return c.sseHub
}

func (c *DIContainer) GetSSEConnector() (*sharedsse.Connector, error) {
	if c.sseConnector == nil {
		connector, err := sharedsse.NewConnector(
			c.GetEventBus(),
			c.GetSSEHub(),
			[]string{avataroptionusecases.AvatarGenerationCompletedEventName},
		)
		if err != nil {
			return nil, err
		}

		c.sseConnector = connector
	}

	return c.sseConnector, nil
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

func (c *DIContainer) GetAvatarOptionRepo(ctx context.Context) (avataroptiondomain.AvatarOptionRepository, error) {
	if c.avatarOptionRepo == nil {
		db, err := c.GetDB(ctx)
		if err != nil {
			return nil, err
		}

		c.avatarOptionRepo = avataroptionrepo.NewAvatarOptionPostgresRepo(db)
	}

	return c.avatarOptionRepo, nil
}

func (c *DIContainer) GetAvatarGenerator() (avataroptiondomain.AvatarGenerator, error) {
	if c.avatarGenerator == nil {
		config, err := c.GetConfig()
		if err != nil {
			return nil, err
		}

		generator, err := avataroptiongenerator.NewOpenAIDALLEGenerator(
			config.OpenAIAPIKey,
			config.OpenAIImageModel,
		)
		if err != nil {
			return nil, err
		}

		c.avatarGenerator = generator
	}

	return c.avatarGenerator, nil
}

func (c *DIContainer) GetFileStorage() (shareddomain.FileStorage, error) {
	if c.fileStorage == nil {
		config, err := c.GetConfig()
		if err != nil {
			return nil, err
		}

		storage, err := sharedstorage.NewLocalFileStorage(
			config.LocalFileStorageRootPath,
			config.BackendPublicBaseURL,
			config.PublicFileURLPrefix,
		)
		if err != nil {
			return nil, err
		}

		c.fileStorage = storage
	}

	return c.fileStorage, nil
}

func (c *DIContainer) GetAvatarConfigRepo(
	ctx context.Context,
) (avatarconfigdomain.AvatarConfigRepository, error) {
	if c.avatarConfigRepo == nil {
		db, err := c.GetDB(ctx)
		if err != nil {
			return nil, err
		}

		c.avatarConfigRepo = avatarconfigrepo.NewAvatarConfigPostgresRepo(db)
	}

	return c.avatarConfigRepo, nil
}
