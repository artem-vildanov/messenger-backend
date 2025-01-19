package routes

import (
	"messenger/internal/bootstrap"
	"messenger/internal/infrastructure/utils/router_utils"
)

func BuildApiGroup(app *bootstrap.App) *router_utils.RoutesGroup {
	authMiddleware := app.MiddlewareRegistry.AuthMiddleware

	authHandler := app.HandlerRegistry.AuthHandler
	userHandler := app.HandlerRegistry.UserHandler
	chatHandler := app.HandlerRegistry.ChatHandler

	return router_utils.NewGroup("/api").
		WithGroups(
			router_utils.NewGroup("/auth").
				WithRoutes(
					router_utils.NewRoute(router_utils.Post, "/login", authHandler.Login),
					router_utils.NewRoute(router_utils.Post, "/register", authHandler.Register),
					router_utils.NewRoute(router_utils.Post, "/logout", authHandler.Logout).Middleware(authMiddleware),
				),
			router_utils.NewGroup("/users").
				WithMiddlewares(authMiddleware).
				WithRoutes(
					router_utils.NewRoute(
						router_utils.Get,
						"/{userId:[0-9]+}",
						userHandler.GetUserById,
					),
					router_utils.NewRoute(
						router_utils.Get,
						"/all",
						userHandler.GetAllUsers,
					),
				),
			router_utils.NewGroup("/chats").
				WithMiddlewares(authMiddleware).
				WithRoutes(
					router_utils.NewRoute(
						router_utils.Get,
						"",
						chatHandler.GetMyChats,
					),
					router_utils.NewRoute(
						router_utils.Get,
						"/{userId:[0-9]+}",
						chatHandler.GetChatMessages,
					),
					// router.Route(router.Post, "/new-message", chatHandler.SendMessage),
				),
		)
}
