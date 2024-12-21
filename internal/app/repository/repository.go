package repo

import (
	"context"
	"messenger/internal/infrastructure/clients"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	postgresClient *clients.PostgresClient
}

func (r *Repository) Construct(postgresClient *clients.PostgresClient) {
	r.postgresClient = postgresClient
}

func (r *Repository) database() *pgxpool.Pool {
	return r.postgresClient.GetClient()
}

func (r *Repository) exec(ctx context.Context, sql string, args ...any) error {
	_, err := r.database().Exec(ctx, sql, args)
	return err
}

func (r *Repository) queryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return r.database().QueryRow(ctx, sql, args)
}

func (r *Repository) query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return r.database().Query(ctx, sql, args)
}
