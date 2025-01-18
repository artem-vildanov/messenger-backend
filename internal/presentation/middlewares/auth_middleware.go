package middlewares

import (
	"messenger/internal/infrastructure/errors"
	"messenger/internal/domain/services"
	"messenger/internal/infrastructure/utils/handler_utils"
	"messenger/internal/infrastructure/utils/router_utils"
)

const failedToAuth = "failed to authenticate user"

func NewAuthMiddleware(sessionService *services.SessionService) router_utils.Middleware {
	return func(
		handlerContext *handler_utils.HandlerContext, 
		next router_utils.Handler,
	) *errors.Error {
		sessionCookie, err := handlerContext.SessionCookie()
		if err != nil {
			return err.WithLogMessage(failedToAuth)
		}

		session, err := sessionService.AuthenticateBySessionId(
			handlerContext.Request.Context(),
			sessionCookie.Value,
		)

		if err != nil {
			return err.WithLogMessage(failedToAuth)
		}

		handlerContext.AuthUserId = session.UserId
		handlerContext.SessionId = session.Id

		return next(handlerContext)
	}
}
