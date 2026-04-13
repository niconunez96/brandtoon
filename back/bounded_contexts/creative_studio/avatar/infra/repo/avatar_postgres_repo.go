package avatarrepo

import (
	avatardomain "brandtoonapi/bounded_contexts/creative_studio/avatar/domain"
	sharedrepos "brandtoonapi/bounded_contexts/shared/infra/repos"
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type AvatarPostgresRepo struct {
	*sharedrepos.PostgresRepo[*avatarDBModel]
	db *sqlx.DB
}

func NewAvatarPostgresRepo(db *sqlx.DB) *AvatarPostgresRepo {
	return &AvatarPostgresRepo{
		PostgresRepo: sharedrepos.NewPostgresRepo(db, func() *avatarDBModel {
			return &avatarDBModel{}
		}),
		db: db,
	}
}

func (r *AvatarPostgresRepo) Create(ctx context.Context, avatar avatardomain.Avatar) error {
	return r.PostgresRepo.Create(ctx, newAvatarDBModel(avatar))
}

func (r *AvatarPostgresRepo) ListByUserID(
	ctx context.Context,
	userID string,
) ([]avatardomain.Avatar, error) {
	models := []*avatarDBModel{}
	err := r.db.SelectContext(
		ctx,
		&models,
		`SELECT * FROM avatars
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
		`SELECT * FROM avatars
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
