package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const defaultServerAddress = "127.0.0.1:8888"

type Config struct {
	AuthStateSecret          string
	BackendPublicBaseURL     string
	DatabaseURL              string
	FrontendBaseURL          string
	GoogleClientID           string
	GoogleClientSecret       string
	GoogleRedirectURL        string
	LocalFileStorageRootPath string
	OpenAIAPIKey             string
	PublicFileURLPrefix      string
	ServerAddress            string
	SessionTTL               time.Duration
}

type UpdateConficFunc func(config *Config, value string)

type EnvConfig struct {
	UpdateConfig UpdateConficFunc
	Required     bool
}

var ENVS = map[string]EnvConfig{
	"DATABASE_URL": {
		Required:     true,
		UpdateConfig: func(config *Config, value string) { config.DatabaseURL = value },
	},
	"BACKEND_PUBLIC_BASE_URL": {
		Required: true,
		UpdateConfig: func(config *Config, value string) {
			config.BackendPublicBaseURL = strings.TrimRight(strings.TrimSpace(value), "/")
		},
	},
	"FRONTEND_BASE_URL": {
		Required:     true,
		UpdateConfig: func(config *Config, value string) { config.FrontendBaseURL = value },
	},
	"GOOGLE_CLIENT_ID": {
		Required:     true,
		UpdateConfig: func(config *Config, value string) { config.GoogleClientID = value },
	},
	"GOOGLE_CLIENT_SECRET": {
		Required:     true,
		UpdateConfig: func(config *Config, value string) { config.GoogleClientSecret = value },
	},
	"GOOGLE_REDIRECT_URL": {
		Required:     true,
		UpdateConfig: func(config *Config, value string) { config.GoogleRedirectURL = value },
	},
	"AUTH_STATE_SECRET": {
		Required:     true,
		UpdateConfig: func(config *Config, value string) { config.AuthStateSecret = value },
	},
	"LOCAL_FILE_STORAGE_ROOT_PATH": {
		Required: false,
		UpdateConfig: func(config *Config, value string) {
			trimmed := strings.TrimSpace(value)
			if trimmed == "" {
				config.LocalFileStorageRootPath = "./storage"
				return
			}

			config.LocalFileStorageRootPath = trimmed
		},
	},
	"OPENAI_API_KEY": {
		Required:     true,
		UpdateConfig: func(config *Config, value string) { config.OpenAIAPIKey = strings.TrimSpace(value) },
	},
	"PUBLIC_FILE_URL_PREFIX": {
		Required: false,
		UpdateConfig: func(config *Config, value string) {
			trimmed := strings.TrimSpace(value)
			if trimmed == "" {
				config.PublicFileURLPrefix = "/files"
				return
			}

			if !strings.HasPrefix(trimmed, "/") {
				trimmed = "/" + trimmed
			}

			config.PublicFileURLPrefix = strings.TrimRight(trimmed, "/")
		},
	},
	"SERVER_ADDRESS": {
		Required: false,
		UpdateConfig: func(config *Config, value string) {
			address := strings.TrimSpace(value)
			if address == "" {
				config.ServerAddress = defaultServerAddress
				return
			}
			config.ServerAddress = address

		},
	},
}

func LoadConfig() (Config, error) {
	config := &Config{
		LocalFileStorageRootPath: "./storage",
		PublicFileURLPrefix:      "/files",
		SessionTTL:               30 * 24 * time.Hour,
	}
	for envKey, envConfig := range ENVS {
		value := os.Getenv(envKey)
		if envConfig.Required {
			requiredValue, err := requiredEnv(envKey)
			if err != nil {
				return Config{}, err
			}
			value = requiredValue
		}
		envConfig.UpdateConfig(config, value)
	}
	return *config, nil
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
