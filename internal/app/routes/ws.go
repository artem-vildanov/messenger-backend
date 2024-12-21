package routes

import (
	"messenger/internal/infrastructure/di"
	"messenger/internal/app/middlewares"
	"messenger/internal/infrastructure/server/router"
)

func Ws(container *di.DependencyContainer) *router.RoutesGroup {
	return router.NewGroup("/ws").WithMiddlewares(
		di.Provide[middlewares.LoggingMiddleware](container),
	).WithGroups(
		router.NewGroup("/messages").WithRoutes(

		),
		router.NewGroup("/chats").WithRoutes(

		),
	)
}