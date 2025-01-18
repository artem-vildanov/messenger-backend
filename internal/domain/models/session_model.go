package models

import (
	"context"
	"messenger/internal/infrastructure/errors"
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

func ExtractSessionFromContext(context context.Context) (*SessionModel, *errors.Error) {
	session, ok := context.Value(handler_utils.ContextKey(handler_utils.SessionIdKey)).(*SessionModel)
	if !ok {
		return nil, errors.InternalError().
			WithLogMessage("session not found in context")
	}

	if session.Id == "" || session.UserId == 0 {
		return nil, errors.InternalError().
			WithLogMessage("invalid session not provided").
			WithField("SessionModel", session)
	}

	return session, nil
}
