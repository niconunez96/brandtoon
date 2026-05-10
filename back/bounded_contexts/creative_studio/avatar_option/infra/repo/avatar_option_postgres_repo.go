package avataroptionrepo

import (
	avataroptiondomain "brandtoonapi/bounded_contexts/creative_studio/avatar_option/domain"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type AvatarOptionPostgresRepo struct{ db *sqlx.DB }

type avatarOptionDBModel struct {
	ID        string     `db:"id"`
	AvatarID  string     `db:"avatar_id"`
	Status    string     `db:"status"`
	Selected  bool       `db:"selected"`
	Href      *string    `db:"href"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func NewAvatarOptionPostgresRepo(db *sqlx.DB) *AvatarOptionPostgresRepo {
	return &AvatarOptionPostgresRepo{db: db}
}

func (r *AvatarOptionPostgresRepo) CreatePendingBatch(
	ctx context.Context,
	options []avataroptiondomain.AvatarOption,
) error {
	now := time.Now().UTC()
	for _, option := range options {
		if _, err := r.db.ExecContext(
			ctx,
			`INSERT INTO avatar_options (
				id, avatar_id, status, selected, href, created_at, updated_at, deleted_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, NULL)`,
			option.ID,
			option.AvatarID,
			string(option.Status),
			option.Selected,
			option.Href,
			now,
			now,
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *AvatarOptionPostgresRepo) ListByAvatarID(
	ctx context.Context,
	avatarID string,
) ([]avataroptiondomain.AvatarOption, error) {
	models := []avatarOptionDBModel{}
	if err := r.db.SelectContext(
		ctx,
		&models,
		`SELECT id, avatar_id, status, selected, href, deleted_at
		 FROM avatar_options
		 WHERE avatar_id = $1 AND deleted_at IS NULL
		 ORDER BY created_at DESC`,
		avatarID,
	); err != nil {
		return nil, err
	}

	options := make([]avataroptiondomain.AvatarOption, 0, len(models))
	for _, model := range models {
		options = append(
			options,
			avataroptiondomain.AvatarOption{
				ID:       model.ID,
				AvatarID: model.AvatarID,
				Status:   avataroptiondomain.Status(model.Status),
				Selected: model.Selected,
				Href:     model.Href,
			},
		)
	}
	return options, nil
}

func (r *AvatarOptionPostgresRepo) Select(ctx context.Context, avatarID string, optionID string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.ExecContext(
		ctx,
		`UPDATE avatar_options
		 SET selected = false, updated_at = $2
		 WHERE avatar_id = $1 AND deleted_at IS NULL`,
		avatarID,
		time.Now().UTC(),
	); err != nil {
		return err
	}

	result, err := tx.ExecContext(
		ctx,
		`UPDATE avatar_options
		 SET selected = true, updated_at = $3
		 WHERE avatar_id = $1 AND id = $2 AND status = 'DONE' AND deleted_at IS NULL`,
		avatarID,
		optionID,
		time.Now().UTC(),
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return avataroptiondomain.ErrAvatarOptionNotFound
	}

	return tx.Commit()
}

func (r *AvatarOptionPostgresRepo) Delete(ctx context.Context, avatarID string, ids []string) error {
	now := time.Now().UTC()
	for _, id := range ids {
		if _, err := r.db.ExecContext(
			ctx,
			`UPDATE avatar_options
			 SET deleted_at = $3, updated_at = $3, selected = false
			 WHERE avatar_id = $1 AND id = $2 AND deleted_at IS NULL`,
			avatarID,
			id,
			now,
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *AvatarOptionPostgresRepo) MarkDone(ctx context.Context, id string, href string) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE avatar_options SET status = 'DONE', href = $2, updated_at = $3 WHERE id = $1 AND deleted_at IS NULL`,
		id,
		href,
		time.Now().UTC(),
	)
	return err
}

func (r *AvatarOptionPostgresRepo) MarkFailed(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE avatar_options
		 SET status = 'FAILED', href = NULL, selected = false, updated_at = $2
		 WHERE id = $1 AND deleted_at IS NULL`,
		id,
		time.Now().UTC(),
	)
	return err
}
