package middlewares

import (
	"log"
	"messenger/internal/app/errors"
	ctx "messenger/internal/infrastructure/handler_context"
)

type LoggingMiddleware struct {
	// there might be a logging or metrics client
}

func (m *LoggingMiddleware) MiddlewareFunc(context *ctx.HandlerContext) *errors.Error {
	log.Printf("Request: %s, %s\n", context.Request.Method, context.Request.URL.Path)
	return nil
}
