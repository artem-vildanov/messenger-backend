package inits

import (
	"fmt"
	"log"
	"messenger/internal/infrastructure/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const pgClientName = "postgres"

func InitPostgres(env *config.Env) (*sqlx.DB, func() error) {
	db, err := sqlx.Connect(pgClientName, fmt.Sprintf(
		`host=%s 
		port=%s 
		user=%s 
		password=%s 
		dbname=%s 
		sslmode=%s`,
		env.PgHost,
		env.PgPort,
		env.PgUser,
		env.PgPassword,
		env.PgDb,
		env.PgSSL,
	))
	if err != nil {
		log.Panicf("failed to connect to postgres: %s", err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Panicf("failed to connect to postgres: %s", err.Error())
	}

	return db, db.Close
}
