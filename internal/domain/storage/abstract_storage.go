package storage

import (
	"context"
	"database/sql"
	appErrors "messenger/internal/infrastructure/errors"

	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// T is type of model with `db:"..."` tags on fields
type AbstractStorage[T any] struct {
	postgres *sqlx.DB
}

func (r *AbstractStorage[T]) exec(ctx context.Context, sql string, args ...any) error {
	_, err := r.postgres.ExecContext(ctx, sql, args...)
	return err
}

func (r *AbstractStorage[T]) queryRow(
	ctx context.Context,
	sql string,
	args ...any,
) *sql.Row {
	return r.postgres.QueryRowContext(
		ctx,
		sql,
		args...,
	)
}

func (r *AbstractStorage[T]) findOne(
	ctx context.Context, 
	query string, 
	args ...any,
) (*T, *appErrors.Error) {
	model := new(T)
	
	if err := r.postgres.GetContext(ctx, model, query, args...); err != nil {
		if r.isNoRows(err) {
			return model, appErrors.NotFoundError()
		}
		return model, appErrors.InternalError().WithLogMessage(err.Error())
	}
	return model, nil
}

func (r *AbstractStorage[T]) findSlice(
	ctx context.Context,
	sql string,
	args ...any,
) ([]*T, *appErrors.Error) {
	models := make([]*T, 0)
	if err := r.postgres.SelectContext(ctx, &models, sql, args...); err != nil {
		return models, appErrors.InternalError().WithLogMessage(err.Error())
	}
	return models, nil
}

func (r *AbstractStorage[T]) isUniqueViolation(err error) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func (r *AbstractStorage[T]) isNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
