package storage

import (
	"context"
	"database/sql"
	appErrors "messenger/internal/infrastructure/errors"
	"errors"
	"github.com/jmoiron/sqlx"
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
) (*T, error) {
	model := new(T)
	
	if err := r.postgres.GetContext(ctx, model, query, args...); err != nil {
		if appErrors.IsNoRowsErr(err) {
			return model, appErrors.Wrap(
				appErrors.ErrNotFound,
				errors.New("findOne"),
			)
		}
		return model, appErrors.Wrap(
			appErrors.ErrInternal, 
			err,
			errors.New("findOne"),
		)
	}
	return model, nil
}

func (r *AbstractStorage[T]) findSlice(
	ctx context.Context,
	sql string,
	args ...any,
) ([]*T, error) {
	models := make([]*T, 0)
	if err := r.postgres.SelectContext(
		ctx, 
		&models, 
		sql, 
		args...,
	); err != nil {
		return models, appErrors.Wrap(
			appErrors.ErrInternal,
			err,
			errors.New("findSlice"),
		)
	}
	return models, nil
}

