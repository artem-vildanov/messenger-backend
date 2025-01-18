package routes

import (
	"messenger/internal/bootstrap"
	"messenger/internal/infrastructure/utils/router_utils"
	"messenger/internal/infrastructure/utils/ws_utils"
)

func BuildWsGroup(app *bootstrap.App) *router_utils.RoutesGroup {
	authMiddleware := app.MiddlewareRegistry.AuthMiddleware
	chatHandler := app.HandlerRegistry.ChatHandler

	return router_utils.NewGroup("/ws").
		WithMiddlewares(authMiddleware).
		WithRoutes(
			router_utils.Route(
				router_utils.Get, 
				"/chat-messages/{userId:[0-9]+}", 
				ws_utils.WsHandlers(
					chatHandler.WriteChatMessages,
					chatHandler.ReadChatMessage,
				),
			),
			router_utils.Route(
				router_utils.Get, 
				"/message-notification", 
				ws_utils.WsHandlers(
					chatHandler.WriteMessageNotifications,
				),
			),
		)
}
