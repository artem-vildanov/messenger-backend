package main

import (
	"messenger/internal/presentation/routes"
	"messenger/internal/bootstrap"
)

func main() {
	app := bootstrap.NewApp()
	defer app.Cleanup()
	rootGroup := routes.BuildRootGroup(app)
	app.Run(rootGroup)
}