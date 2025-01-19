package models

import (
	"context"
	"errors"
	appErrors "messenger/internal/infrastructure/errors"
	"messenger/internal/infrastructure/utils/handler_utils"
	"time"
)

type AuthModel struct {
	Username string
	Password string
}

type SessionModel struct {
	Id        string    `db:"id"`
	ExpiresAt time.Time `db:"expires_at"`
	UserId    int       `db:"user_id"`
}

func ExtractSessionFromContext(ctx context.Context) (*SessionModel, error) {
	session, ok := ctx.Value(
		handler_utils.ContextKey(handler_utils.SessionIdKey),
	).(*SessionModel)
	if !ok {
		return nil, appErrors.Wrap(
			appErrors.ErrInternal,
			errors.New("ExtractSessionFromContext"),
		)
	}

	if session.Id == "" || session.UserId == 0 {
		return nil, appErrors.Wrap(
			appErrors.ErrInternal,
			errors.New("ExtractSessionFromContext"),
		)
	}

	return session, nil
}
