package http

import (
	"net/http"
	"time"

	sharedconfig "brandtoonapi/bounded_contexts/shared/infra/config"
)

const sessionCookieName = "brandtoon_session_id"

func newSessionCookie(config sharedconfig.Config, value string, expiresAt time.Time) http.Cookie {
	return http.Cookie{
		Name:     sessionCookieName,
		Value:    value,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
		SameSite: http.SameSiteLaxMode,
		Secure:   config.SessionCookieSecure(),
	}
}

func expiredSessionCookie(config sharedconfig.Config) http.Cookie {
	return http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		Secure:   config.SessionCookieSecure(),
	}
}
