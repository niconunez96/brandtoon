package domain

import "time"

type Session struct {
	ExpiresAt time.Time
	ID        string
	UserID    string
}

func NewSession(id string, userID string, expiresAt time.Time) Session {
	return Session{
		ExpiresAt: expiresAt,
		ID:        id,
		UserID:    userID,
	}
}

func (s Session) IsExpired(now time.Time) bool {
	return !s.ExpiresAt.After(now)
}
