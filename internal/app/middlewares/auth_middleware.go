package middlewares

import (
	"context"
	"messenger/internal/app/errors"
	"messenger/internal/app/handlers/ctx"
	"messenger/internal/app/repository"
)

type AuthMiddleware struct {
	authRepo repo.AuthRepository
}

func (m *AuthMiddleware) Construct(authRepo repo.AuthRepository) {
	m.authRepo = authRepo
}

func (m *AuthMiddleware) MiddlewareFunc(handlerContext *ctx.HandlerContext) *errors.Error {
	requestContext := handlerContext.Request.Context()
	sessionCookie, err := handlerContext.SessionCookie()
	if err != nil {
		return err
	}

	session, err := m.authRepo.GetSession(requestContext, sessionCookie.Value)
	if err != nil {
		return err
	}

	requestContext = context.WithValue(
		requestContext,
		ctx.ContextKey(ctx.SessionKey),
		session,
	)

	handlerContext.Request = handlerContext.Request.WithContext(requestContext)
	return nil
}
