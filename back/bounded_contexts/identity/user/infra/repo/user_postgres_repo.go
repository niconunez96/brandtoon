package repo

import (
	"context"
	"database/sql"
	"errors"

	"brandtoonapi/bounded_contexts/identity/user/domain"
	sharedrepos "brandtoonapi/bounded_contexts/shared/infra/repos"

	"github.com/jmoiron/sqlx"
)

type UserPostgresRepo struct {
	*sharedrepos.PostgresRepo[*userDBModel]
	db *sqlx.DB
}

func NewUserPostgresRepo(db *sqlx.DB) *UserPostgresRepo {
	return &UserPostgresRepo{
		PostgresRepo: sharedrepos.NewPostgresRepo(db, func() *userDBModel {
			return &userDBModel{}
		}),
		db: db,
	}
}

func (r *UserPostgresRepo) Create(ctx context.Context, user domain.User) error {
	return r.PostgresRepo.Create(ctx, newUserDBModel(user))
}

func (r *UserPostgresRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	model := &userDBModel{}
	err := r.db.GetContext(
		ctx,
		model,
		"SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL LIMIT 1",
		email,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	user := model.ToDomain()
	return &user, nil
}

func (r *UserPostgresRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	model, err := r.PostgresRepo.FindByID(ctx, id)
	if err != nil || model == nil {
		return nil, err
	}

	user := model.ToDomain()
	return &user, nil
}

func (r *UserPostgresRepo) Update(ctx context.Context, user domain.User) error {
	return r.PostgresRepo.Update(ctx, newUserDBModel(user))
}
