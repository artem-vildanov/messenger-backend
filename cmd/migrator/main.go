package main

import (
	"log"
	"messenger/internal/infrastructure/config"
	"messenger/internal/infrastructure/inits"

	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	env := config.LoadEnv()
	pg, close := inits.InitPostgres(env)
	defer close()

	amount, err := migrate.Exec(
		pg.DB,
		"postgres",
		&migrate.FileMigrationSource{Dir: "./migrations"},
		migrate.Up,
	)

	if err != nil {
		log.Panic(err)
	}

	log.Printf("Applied %d migrations\n", amount)
}
