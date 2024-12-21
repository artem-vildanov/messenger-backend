package routes

import (
	"messenger/internal/infrastructure/di"
	"messenger/internal/app/handlers"
	"messenger/internal/app/middlewares"
	repo "messenger/internal/app/repository"
	"messenger/internal/infrastructure/server/router"
)

func Api(container *di.DependencyContainer) *router.RoutesGroup {
	di.Bind[repo.AuthRepository, repo.AuthRepositoryImpl](container)
	di.Bind[repo.UserRepository, repo.UserRepositoryImpl](container)

	loggingMiddleware := di.Provide[middlewares.LoggingMiddleware](container)
	authMiddleware := di.Provide[middlewares.AuthMiddleware](container)

	authHandler := di.Provide[handlers.AuthHandler](container)
	userHandler := di.Provide[handlers.UserHandler](container)

	return router.NewGroup("/api").
		WithMiddlewares(loggingMiddleware).
		WithGroups(
			router.NewGroup("/auth").
				WithRoutes(
					router.Route(router.Post, "/login", authHandler.Login),
					router.Route(router.Post, "/register", authHandler.Register),
					router.Route(router.Post, "/logout", authHandler.Logout).Middleware(authMiddleware),
				),
			router.NewGroup("/users").
				WithMiddlewares(authMiddleware).
				WithRoutes(
					router.Route(router.Get, "/{userId:[0-9]+}", userHandler.GetUserById),
				),
		)
}
