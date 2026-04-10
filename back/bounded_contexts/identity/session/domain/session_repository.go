package domain

import "context"

type SessionRepository interface {
	Create(ctx context.Context, session Session) error
	Delete(ctx context.Context, id string) error
	FindActiveByID(ctx context.Context, id string) (*Session, error)
}
