package main

import (
	"context"
	"log"
	"net/http"
	"time"

	avatarhttp "brandtoonapi/bounded_contexts/creative_studio/avatar/infra/http"
	avatarconfighttp "brandtoonapi/bounded_contexts/creative_studio/avatar_config/infra/http"
	authhttp "brandtoonapi/bounded_contexts/identity/auth/infra/http"
	shared "brandtoonapi/bounded_contexts/shared"
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	sharedconfig "brandtoonapi/bounded_contexts/shared/infra/config"
	"brandtoonapi/bounded_contexts/shared/infra/telemetry"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func main() {
	ctx := context.Background()
	_ = godotenv.Load()
	container := shared.NewDIContainer()
	t := telemetry.NewTelemetry()

	config, err := container.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	googleProvider, err := container.GetGoogleIdentityProvider()
	if err != nil {
		log.Fatal(err)
	}

	stateCodec, err := container.GetOAuthStateCodec()
	if err != nil {
		log.Fatal(err)
	}

	userRepo, err := container.GetUserRepo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	sessionRepo, err := container.GetSessionRepo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	avatarRepo, err := container.GetAvatarRepo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	avatarConfigRepo, err := container.GetAvatarConfigRepo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()
	router.Use(corsMiddleware(config))
	router.Use(t.HttpLoggerMiddleware())
	api := humachi.New(router, huma.DefaultConfig("Brandtoon API", "1.0.0"))

	authhttp.RegisterRoutes(api, router, authhttp.RouteDependencies{
		Config:         config,
		GoogleProvider: googleProvider,
		IDGenerator:    shareddomain.GenerateUUIDv7,
		Now:            func() time.Time { return time.Now().UTC() },
		SessionRepo:    sessionRepo,
		StateCodec:     stateCodec,
		UserRepo:       userRepo,
	})
	avatarhttp.RegisterRoutes(api, router, avatarhttp.RouteDependencies{
		AvatarRepo:  avatarRepo,
		IDGenerator: shareddomain.GenerateUUIDv7,
		SessionRepo: sessionRepo,
		UserRepo:    userRepo,
	})
	avatarconfighttp.RegisterRoutes(api, router, avatarconfighttp.RouteDependencies{
		AvatarConfigRepo: avatarConfigRepo,
		AvatarRepo:       avatarRepo,
		SessionRepo:      sessionRepo,
		UserRepo:         userRepo,
	})

	if err := http.ListenAndServe(config.ServerAddress, router); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(config sharedconfig.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			origin := request.Header.Get("Origin")
			if origin == config.FrontendBaseURL {
				writer.Header().Set("Access-Control-Allow-Credentials", "true")
				writer.Header().Set("Access-Control-Allow-Origin", origin)
				writer.Header().Set("Vary", "Origin")
			}

			writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
			writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if request.Method == http.MethodOptions {
				writer.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}
