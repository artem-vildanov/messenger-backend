package clients

import (
	"database/sql"
	"messenger/internal/app/errors"
	"messenger/internal/infrastructure/config"
)

type MigratorClient struct {
	connection *sql.DB
}

func (c *MigratorClient) Construct(env *config.Env) {
	db, err := sql.Open(postgresClientName, env.GetPostgresAddr().WithSllDisabled().String())
	if err != nil {
		errors.ClientConnectionPanic(postgresClientName, err.Error())
	}

	err = db.Ping()
	if err != nil {
		errors.ClientConnectionPanic(postgresClientName, err.Error())
	}

	c.connection = db
}

func (c *MigratorClient) GetClient() *sql.DB {
	return c.connection
}

func (c *MigratorClient) CloseConnection() {
	c.connection.Close()
}
