package postgres

import (
	"context"
	"embed"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func NewDatabase(ctx context.Context, databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, err
	}

	goose.SetBaseFS(migrationsFS)
	if err := goose.UpContext(ctx, db.DB, "migrations"); err != nil {
		return nil, err
	}

	return db, nil
}
