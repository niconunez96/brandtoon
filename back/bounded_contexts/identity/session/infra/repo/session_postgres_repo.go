package sessionrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"brandtoonapi/bounded_contexts/identity/session/domain"
	"brandtoonapi/bounded_contexts/shared/infra/repos"

	"github.com/jmoiron/sqlx"
)

type SessionPostgresRepo struct {
	*sharedrepos.PostgresRepo[*sessionDBModel]
	db *sqlx.DB
}

func NewSessionPostgresRepo(db *sqlx.DB) *SessionPostgresRepo {
	return &SessionPostgresRepo{
		PostgresRepo: sharedrepos.NewPostgresRepo(db, func() *sessionDBModel {
			return &sessionDBModel{}
		}),
		db: db,
	}
}

func (r *SessionPostgresRepo) Create(ctx context.Context, session sessiondomain.Session) error {
	return r.PostgresRepo.Create(ctx, newSessionDBModel(session))
}

func (r *SessionPostgresRepo) Delete(ctx context.Context, id string) error {
	return r.PostgresRepo.Delete(ctx, id)
}

func (r *SessionPostgresRepo) FindActiveByID(ctx context.Context, id string) (*sessiondomain.Session, error) {
	model := &sessionDBModel{}
	err := r.db.GetContext(
		ctx,
		model,
		"SELECT * FROM sessions WHERE id = $1 AND deleted_at IS NULL AND expires_at > $2 LIMIT 1",
		id,
		time.Now().UTC(),
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	session := model.ToDomain()
	return &session, nil
}
