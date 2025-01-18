package middlewares

import (
	"log"
	"messenger/internal/infrastructure/errors"
	ctx "messenger/internal/infrastructure/utils/handler_utils"
	"messenger/internal/infrastructure/utils/router_utils"
)

func NewLogginMiddleware() router_utils.Middleware {
	return func(
		handlerContext *ctx.HandlerContext,
		next router_utils.Handler,
	) *errors.Error {
		log.Printf(
			"Request: %s, %s\n",
			handlerContext.Request.Method,
			handlerContext.Request.URL.Path,
		)
		if err := next(handlerContext); err != nil {
			err.LogStdout()
			return err
		}
		return nil
	}
}
