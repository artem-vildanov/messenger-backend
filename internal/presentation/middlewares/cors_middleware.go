package middlewares

import (
	"messenger/internal/infrastructure/errors"
	"messenger/internal/infrastructure/config"
	ctx "messenger/internal/infrastructure/utils/handler_utils"
	"messenger/internal/infrastructure/utils/router_utils"
)

func NewCorsMiddleware(env *config.Env) router_utils.Middleware {
	return func(
		handlerContext *ctx.HandlerContext,
		next router_utils.Handler,
	) *errors.Error {
		builder := handlerContext.Response().
			WithHeader("Access-Control-Allow-Origin", env.Origin).
			WithHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS").
			WithHeader("Access-Control-Allow-Credentials", "true").
			WithHeader("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if handlerContext.Request.Method == "OPTIONS" {
			builder.Empty()
			return nil
		}

		return next(handlerContext)
	}
}
