package usecases

import "strings"

const defaultRedirectTarget = "/creative-studio"

func sanitizeRedirectTarget(redirectTo string) string {
	trimmed := strings.TrimSpace(redirectTo)
	if trimmed == "" {
		return defaultRedirectTarget
	}

	if !strings.HasPrefix(trimmed, "/") || strings.HasPrefix(trimmed, "//") {
		return defaultRedirectTarget
	}

	return trimmed
}
