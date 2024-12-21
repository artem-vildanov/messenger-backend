package main

import (
	"messenger/internal/infrastructure/server"
)

func main() {
	srv := server.New()
	defer srv.BeforeShutdown()
	srv.Run()
}