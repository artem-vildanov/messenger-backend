package main

import (
	"log"
	"messenger/internal/infrastructure/clients"
	"messenger/internal/infrastructure/di"

	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	container := di.NewDependencyContainer()
	migratorClient := di.Provide[clients.MigratorClient](container)
	defer migratorClient.CloseConnection()

	amount, err := migrate.Exec(
		migratorClient.GetClient(),
		"postgres",
		&migrate.FileMigrationSource{Dir: "./migrations"},
		migrate.Up,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Applied %d migrations\n", amount)
}
