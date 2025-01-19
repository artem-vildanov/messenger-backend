package middlewares

import (
	"errors"
	"messenger/internal/domain/services"
	appErrors "messenger/internal/infrastructure/errors"
	"messenger/internal/infrastructure/utils/handler_utils"
	"messenger/internal/infrastructure/utils/router_utils"
)

const failedToAuth = "failed to authenticate user"

func NewAuthMiddleware(sessionService *services.SessionService) router_utils.Middleware {
	return func(
		handlerContext *handler_utils.HandlerContext, 
		next router_utils.Handler,
	) error {
		sessionCookie, err := handlerContext.SessionCookie()
		if err != nil {
			return appErrors.Wrap(err, errors.New("AuthMiddleware"))
		}

		session, err := sessionService.AuthenticateBySessionId(
			handlerContext.Request.Context(),
			sessionCookie.Value,
		)

		if err != nil {
			return appErrors.Wrap(err, errors.New("AuthMiddleware"))
		}

		handlerContext.AuthUserId = session.UserId
		handlerContext.SessionId = session.Id

		return next(handlerContext)
	}
}
