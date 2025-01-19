package middlewares

import (
	"log"
	ctx "messenger/internal/infrastructure/utils/handler_utils"
	"messenger/internal/infrastructure/utils/router_utils"
)

func NewLogginMiddleware() router_utils.Middleware {
	return func(
		handlerContext *ctx.HandlerContext,
		next router_utils.Handler,
	) error {
		log.Printf(
			"Request: %s, %s\n",
			handlerContext.Request.Method,
			handlerContext.Request.URL.Path,
		)
		return next(handlerContext)
	}
}
