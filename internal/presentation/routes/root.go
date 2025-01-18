package routes

import (
	"messenger/internal/bootstrap"
	"messenger/internal/infrastructure/utils/router_utils"
)

func BuildRootGroup(app *bootstrap.App) *router_utils.RoutesGroup {
	corsMiddleware := app.MiddlewareRegistry.CorsMiddleware
	loggingMiddleware := app.MiddlewareRegistry.LoggingMiddleware

	return router_utils.RootGroup().
		WithMiddlewares(corsMiddleware, loggingMiddleware).
		WithGroups(
			BuildApiGroup(app),
			BuildWsGroup(app),
		)
}
