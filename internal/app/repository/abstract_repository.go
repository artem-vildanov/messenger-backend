package repository

import (
	"context"
	"messenger/internal/infrastructure/clients"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AbstractRepository struct {
	postgresClient *clients.PostgresClient
}

func (r *AbstractRepository) Construct(postgresClient *clients.PostgresClient) {
	r.postgresClient = postgresClient
}

func (r *AbstractRepository) database() *pgxpool.Pool {
	return r.postgresClient.GetClient()
}

func (r *AbstractRepository) exec(ctx context.Context, sql string, args ...any) error {
	_, err := r.database().Exec(ctx, sql, args...)
	return err
}

func (r *AbstractRepository) queryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return r.database().QueryRow(ctx, sql, args...)
}

func (r *AbstractRepository) query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return r.database().Query(ctx, sql, args...)
}
