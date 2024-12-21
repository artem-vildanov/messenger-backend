package middlewares

import (
	"messenger/internal/app/errors"
	"messenger/internal/app/handlers/ctx"
)

type Middleware interface {
	MiddlewareFunc(*ctx.HandlerContext) *errors.Error
}
