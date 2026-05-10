package avatarrepo

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type AvatarPostgresRepo struct {
	db *sqlx.DB
}

func NewAvatarPostgresRepo(db *sqlx.DB) *AvatarPostgresRepo {
	return &AvatarPostgresRepo{db: db}
}

func (r *AvatarPostgresRepo) Create(ctx context.Context, avatar avatardomain.Avatar) error {
	now := time.Now().UTC()
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO avatars (id, user_id, name, created_at, updated_at, deleted_at)
		 VALUES ($1, $2, $3, $4, $5, NULL)`,
		avatar.ID,
		avatar.UserID,
		avatar.Name,
		now,
		now,
	)
	return err
}

func (r *AvatarPostgresRepo) ListByUserID(
	ctx context.Context,
	userID string,
) ([]avatardomain.Avatar, error) {
	models := []*avatarDBModel{}
	err := r.db.SelectContext(
		ctx,
		&models,
		`SELECT id, user_id, name, created_at, updated_at, deleted_at
		 FROM avatars
		 WHERE user_id = $1 AND deleted_at IS NULL
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	avatars := make([]avatardomain.Avatar, 0, len(models))
	for _, model := range models {
		avatars = append(avatars, model.ToDomain())
	}

	return avatars, nil
}

func (r *AvatarPostgresRepo) FindOwnedByID(
	ctx context.Context,
	avatarID string,
	userID string,
) (*avatardomain.Avatar, error) {
	model := &avatarDBModel{}
	err := r.db.GetContext(
		ctx,
		model,
		`SELECT id, user_id, name, created_at, updated_at, deleted_at
		 FROM avatars
		 WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
		 LIMIT 1`,
		avatarID,
		userID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	avatar := model.ToDomain()
	return &avatar, nil
}
