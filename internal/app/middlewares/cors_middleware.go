package middlewares

import (
	"messenger/internal/app/errors"
	"messenger/internal/infrastructure/config"
	ctx "messenger/internal/infrastructure/handler_context"
)

type CorsMiddleware struct {
	env *config.Env
}

func (m *CorsMiddleware) Construct(env *config.Env) {
	m.env = env
}

func (m *CorsMiddleware) MiddlewareFunc(handlerContext *ctx.HandlerContext) *errors.Error {
	if handlerContext.Request.Method == "OPTIONS" {
		handlerContext.Response().
			WithHeader("Access-Control-Allow-Origin", m.env.GetOrigin()).
			WithHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS").
			WithHeader("Access-Control-Allow-Credentials", "true").
			WithHeader("Access-Control-Allow-Headers", "*").
			Empty()
		return nil
	}

	handlerContext.Response().WithHeader("Access-Control-Allow-Origin", m.env.GetOrigin())
	return nil
}
