package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const defaultServerAddress = "127.0.0.1:8888"

type Config struct {
	AuthStateSecret    string
	DatabaseURL        string
	FrontendBaseURL    string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	ServerAddress      string
	SessionTTL         time.Duration
}

func LoadConfig() (Config, error) {
	databaseURL, err := requiredEnv("DATABASE_URL")
	if err != nil {
		return Config{}, err
	}

	frontendBaseURL, err := requiredEnv("FRONTEND_BASE_URL")
	if err != nil {
		return Config{}, err
	}

	googleClientID, err := requiredEnv("GOOGLE_CLIENT_ID")
	if err != nil {
		return Config{}, err
	}

	googleClientSecret, err := requiredEnv("GOOGLE_CLIENT_SECRET")
	if err != nil {
		return Config{}, err
	}

	googleRedirectURL, err := requiredEnv("GOOGLE_REDIRECT_URL")
	if err != nil {
		return Config{}, err
	}

	authStateSecret, err := requiredEnv("AUTH_STATE_SECRET")
	if err != nil {
		return Config{}, err
	}

	serverAddress := strings.TrimSpace(os.Getenv("SERVER_ADDRESS"))
	if serverAddress == "" {
		serverAddress = defaultServerAddress
	}

	return Config{
		AuthStateSecret:    authStateSecret,
		DatabaseURL:        databaseURL,
		FrontendBaseURL:    strings.TrimRight(frontendBaseURL, "/"),
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		GoogleRedirectURL:  googleRedirectURL,
		ServerAddress:      serverAddress,
		SessionTTL:         30 * 24 * time.Hour,
	}, nil
}

func (c Config) SessionCookieSecure() bool {
	return strings.HasPrefix(c.FrontendBaseURL, "https://")
}

func requiredEnv(key string) (string, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return "", fmt.Errorf("missing required env var %s", key)
	}

	return value, nil
}
