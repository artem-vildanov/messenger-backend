package routes

import (
	"messenger/internal/bootstrap"
	utils "messenger/internal/infrastructure/utils/router_utils"
)

func BuildApiGroup(app *bootstrap.App) *utils.RoutesGroup {
	authMiddleware := app.MiddlewareRegistry.AuthMiddleware

	authHandler := app.HandlerRegistry.AuthHandler
	userHandler := app.HandlerRegistry.UserHandler
	chatHandler := app.HandlerRegistry.ChatHandler

	return utils.NewGroup("/api").
		WithGroups(
			utils.NewGroup("/auth").
				WithRoutes(
					utils.NewRoute(utils.Get, "", userHandler.GetMyUser).Middleware(authMiddleware),
					utils.NewRoute(utils.Post, "/login", authHandler.Login),
					utils.NewRoute(utils.Post, "/register", authHandler.Register),
					utils.NewRoute(utils.Post, "/logout", authHandler.Logout).Middleware(authMiddleware),
				),
			utils.NewGroup("/users").
				WithMiddlewares(authMiddleware).
				WithRoutes(
					utils.NewRoute(utils.Get, "/{userId:[0-9]+}", userHandler.GetUserById),
					utils.NewRoute(utils.Get, "/all", userHandler.GetAllUsers),
				),
			utils.NewGroup("/chats").
				WithMiddlewares(authMiddleware).
				WithRoutes(
					utils.NewRoute(utils.Get, "", chatHandler.GetMyChats),
					utils.NewRoute(utils.Get, "/{userId:[0-9]+}", chatHandler.GetChatMessages),
					utils.NewRoute(
						utils.Get,
						"/{userId:[0-9]+}/missed-messages/{lastMessageId:[0-9]+}",
						chatHandler.GetMissedMessages,
					),
					// router.Route(router.Post, "/new-message", chatHandler.SendMessage),
				),
		)
}
