package avatarconfigrepo

import (
	avatarconfigdomain "brandtoonapi/bounded_contexts/creative_studio/avatar_config/domain"
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type AvatarConfigPostgresRepo struct {
	db *sqlx.DB
}

func NewAvatarConfigPostgresRepo(db *sqlx.DB) *AvatarConfigPostgresRepo {
	return &AvatarConfigPostgresRepo{db: db}
}

func (r *AvatarConfigPostgresRepo) FindByAvatarID(
	ctx context.Context,
	avatarID string,
) (*avatarconfigdomain.AvatarConfig, error) {
	model := &avatarConfigDBModel{}
	err := r.db.GetContext(
		ctx,
		model,
		`SELECT avatar_id, prompt, artistic_style
		 FROM avatar_configs
		 WHERE avatar_id = $1
		 LIMIT 1`,
		avatarID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	avatarConfig, err := model.ToDomain()
	if err != nil {
		return nil, err
	}

	return &avatarConfig, nil
}

func (r *AvatarConfigPostgresRepo) Upsert(
	ctx context.Context,
	avatarConfig avatarconfigdomain.AvatarConfig,
) error {
	now := time.Now().UTC()
	model := newAvatarConfigDBModel(avatarConfig)
	_, err := r.db.NamedExecContext(
		ctx,
		`INSERT INTO avatar_configs (avatar_id, prompt, artistic_style, created_at, updated_at)
		 VALUES (:avatar_id, :prompt, :artistic_style, :created_at, :updated_at)
		 ON CONFLICT (avatar_id) DO UPDATE SET
		   prompt = EXCLUDED.prompt,
		   artistic_style = EXCLUDED.artistic_style,
		   updated_at = EXCLUDED.updated_at`,
		map[string]any{
			"avatar_id":      model.AvatarID,
			"artistic_style": model.ArtisticStyle,
			"created_at":     now,
			"prompt":         model.Prompt,
			"updated_at":     now,
		},
	)
	return err
}
