package clients

import (
	"context"
	"messenger/internal/infrastructure/config"
	"messenger/internal/app/errors"

	pg "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

const postgresClientName = "postgres"

type PostgresClient struct {
	connection *pg.Pool
}

func (d *PostgresClient) Construct(env *config.Env) {
	var err error

	d.connection, err = pg.New(context.Background(), env.GetPostgresAddr().String())
	if err != nil {
		errors.ClientConnectionPanic(postgresClientName, err.Error())
	}

	if err = d.connection.Ping(context.Background()); err != nil {
		errors.ClientConnectionPanic(postgresClientName, err.Error())
	}
}

func (d *PostgresClient) CloseConnection() {
	d.connection.Close()
}

func (d *PostgresClient) GetClient() *pg.Pool {
	return d.connection
}
