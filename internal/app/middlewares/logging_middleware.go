package middlewares

import (
	"log"
	"messenger/internal/app/errors"
	"messenger/internal/app/handlers/ctx"
)

type LoggingMiddleware struct {
	// there might be a logging or metrics client
}

func (m *LoggingMiddleware) MiddlewareFunc(context *ctx.HandlerContext) *errors.Error {
	log.Printf("Request: %s, %s\n", context.Request.Method, context.Request.URL.Path)
	return nil
}
