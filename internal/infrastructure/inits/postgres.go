package inits

import (
	"fmt"
	"log"
	"messenger/internal/infrastructure/config"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitPostgres(env *config.Env) (*sqlx.DB, func() error) {
	var (
		db  *sqlx.DB
		err error
	)

	for retry := range env.PgConnectRetries {
		db, err = sqlx.Connect("postgres", fmt.Sprintf(
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
		if err == nil {
			break
		}
		log.Printf("pg connect retry %d...", retry)
		time.Sleep(time.Second)
	}

	if err != nil {
		log.Panicf("failed to connect to postgres: %s", err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Panicf("failed to connect to postgres: %s", err.Error())
	}

	return db, db.Close
}
