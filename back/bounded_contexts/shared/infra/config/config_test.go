package config

import "testing"

func TestLoadConfigLoadsAvatarGenerationSettings(t *testing.T) {
	t.Setenv("AUTH_STATE_SECRET", "secret")
	t.Setenv("BACKEND_PUBLIC_BASE_URL", "http://127.0.0.1:8888/")
	t.Setenv("DATABASE_URL", "postgres://brandtoon")
	t.Setenv("FRONTEND_BASE_URL", "http://localhost:5173")
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	t.Setenv("GOOGLE_REDIRECT_URL", "http://localhost:8888/auth/google/callback")
	t.Setenv("OPENAI_API_KEY", "sk-test")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if config.BackendPublicBaseURL != "http://127.0.0.1:8888" {
		t.Fatalf("unexpected backend public base url %q", config.BackendPublicBaseURL)
	}

	if config.LocalFileStorageRootPath != "./storage" {
		t.Fatalf("unexpected storage root %q", config.LocalFileStorageRootPath)
	}

	if config.PublicFileURLPrefix != "/files" {
		t.Fatalf("unexpected public file prefix %q", config.PublicFileURLPrefix)
	}
}

