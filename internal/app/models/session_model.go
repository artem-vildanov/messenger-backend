package models

import (
	"context"
	"messenger/internal/app/errors"
	"messenger/internal/app/handlers/ctx"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SessionModel struct {
	ID        string
	UserId    int
	ExpiredAt string
}

func NewSession(userId int, sessionTTL int) *SessionModel {
	return &SessionModel{
		ID: uuid.NewString(),
		UserId: userId,
		ExpiredAt: time.Now().Add(
			time.Duration(sessionTTL)*time.Minute,
		).Format(time.RFC3339),
	}
}

func (m *SessionModel) FromDb(row pgx.Row) *errors.SessionError {
	if err := row.Scan(
		&m.ID,
		&m.UserId,
		&m.ExpiredAt,
	); err != nil {
		if pgx.ErrNoRows.Error() == err.Error() {
			return errors.SessionNotFoundError()
		}
		return errors.FailedToFindSession().WithReason(err.Error())
	}
	return nil
}

func (m *SessionModel) FromContext(context context.Context) *errors.Error {
	session, ok := context.Value(ctx.ContextKey(ctx.SessionKey)).(*SessionModel)
	if !ok {
		return errors.InternalError().WithVerbose("session not found in context")
	}

	*m = *session
	return nil
}
