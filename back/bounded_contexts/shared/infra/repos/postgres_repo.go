package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type PersistableDBModel interface {
	GetID() string
	InsertValues() map[string]any
	SetCreatedAt(time.Time)
	SetDeletedAt(*time.Time)
	SetUpdatedAt(time.Time)
	TableName() string
	UpdateValues() map[string]any
}

type PostgresRepo[T PersistableDBModel] struct {
	db       *sqlx.DB
	newModel func() T
}

func NewPostgresRepo[T PersistableDBModel](db *sqlx.DB, newModel func() T) *PostgresRepo[T] {
	return &PostgresRepo[T]{
		db:       db,
		newModel: newModel,
	}
}

func (r *PostgresRepo[T]) Create(ctx context.Context, model T) error {
	now := time.Now().UTC()
	model.SetCreatedAt(now)
	model.SetUpdatedAt(now)
	model.SetDeletedAt(nil)

	query, values := buildInsertQuery(model.TableName(), model.InsertValues())
	_, err := r.db.NamedExecContext(ctx, query, values)
	return err
}

func (r *PostgresRepo[T]) Update(ctx context.Context, model T) error {
	now := time.Now().UTC()
	model.SetUpdatedAt(now)

	query, values := buildUpdateQuery(model.TableName(), model.UpdateValues(), model.GetID())
	_, err := r.db.NamedExecContext(ctx, query, values)
	return err
}

func (r *PostgresRepo[T]) Delete(ctx context.Context, id string) error {
	model := r.newModel()
	now := time.Now().UTC()
	query := fmt.Sprintf(
		"UPDATE %s SET deleted_at = :deleted_at, updated_at = :updated_at WHERE id = :id AND deleted_at IS NULL",
		model.TableName(),
	)

	_, err := r.db.NamedExecContext(ctx, query, map[string]any{
		"deleted_at": now,
		"id":         id,
		"updated_at": now,
	})
	return err
}

func (r *PostgresRepo[T]) FindByID(ctx context.Context, id string) (T, error) {
	model := r.newModel()
	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE id = $1 AND deleted_at IS NULL LIMIT 1",
		model.TableName(),
	)

	err := r.db.GetContext(ctx, model, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		var zero T
		return zero, nil
	}

	if err != nil {
		var zero T
		return zero, err
	}

	return model, nil
}

func buildInsertQuery(tableName string, values map[string]any) (string, map[string]any) {
	columns := sortedKeys(values)
	placeholders := make([]string, 0, len(columns))
	for _, column := range columns {
		placeholders = append(placeholders, fmt.Sprintf(":%s", column))
	}

	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	), values
}

func buildUpdateQuery(tableName string, values map[string]any, id string) (string, map[string]any) {
	valuesWithID := maps.Clone(values)
	valuesWithID["id"] = id

	columns := sortedKeys(values)
	assignments := make([]string, 0, len(columns))
	for _, column := range columns {
		assignments = append(assignments, fmt.Sprintf("%s = :%s", column, column))
	}

	return fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = :id AND deleted_at IS NULL",
		tableName,
		strings.Join(assignments, ", "),
	), valuesWithID
}

func sortedKeys(values map[string]any) []string {
	keys := slices.Collect(maps.Keys(values))
	slices.Sort(keys)
	return keys
}
