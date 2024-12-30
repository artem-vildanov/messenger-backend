package middlewares

import (
	"context"
	"messenger/internal/app/errors"
	"messenger/internal/app/repository"
	ctx "messenger/internal/infrastructure/handler_context"
	"messenger/internal/infrastructure/server/router"
)

type AuthMiddleware struct {
	authRepo repository.SessionRepository
}

func (m *AuthMiddleware) Construct(authRepo repository.SessionRepository) {
	m.authRepo = authRepo
}

func (m *AuthMiddleware) MiddlewareFunc(
	handlerContext *ctx.HandlerContext, 
	next router.HandlerFunction,
) *errors.Error {
	requestContext := handlerContext.Request.Context()
	sessionCookie, err := handlerContext.SessionCookie()
	if err != nil {
		return err
	}

	session, err := m.authRepo.GetSession(requestContext, sessionCookie.Value)
	if err != nil {
		return err
	}

	if err := session.CheckExpired(); err != nil {
		if err := m.authRepo.DeleteSession(requestContext, session.ID); err != nil {
			return err
		}

		return err.WithId(session.ID).
			WithUserId(session.UserId).
			BuildError()
	}

	requestContext = context.WithValue(
		requestContext,
		ctx.ContextKey(ctx.SessionKey),
		session,
	)

	handlerContext.Request = handlerContext.Request.WithContext(requestContext)

	return next(handlerContext)
}
