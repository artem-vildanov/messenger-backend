package handler_context

import (
	"context"
	"messenger/internal/app/errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const SessionIdKey string = "sessionId"
const SessionKey string = "session"

type Session struct {
	ID        string
	ExpiredAt string
	UserId    int
}

func NewSession(userId int, sessionTTL int) *Session {
	return &Session{
		ID:     uuid.NewString(),
		UserId: userId,
		ExpiredAt: time.Now().Add(
			time.Duration(sessionTTL) * time.Minute,
		).Format(time.RFC3339),
	}
}

func (m *Session) FromDb(row pgx.Row) *errors.SessionError {
	if err := row.Scan(
		&m.ID,
		&m.ExpiredAt,
		&m.UserId,
	); err != nil {
		if pgx.ErrNoRows.Error() == err.Error() {
			return errors.SessionNotFoundError()
		}
		return errors.FailedToFindSession().WithReason(err.Error())
	}
	return nil
}

func (m *Session) FromContext(context context.Context) *errors.Error {
	session, ok := context.Value(ContextKey(SessionKey)).(*Session)
	if !ok {
		return errors.InternalError().WithVerbose("session not found in context")
	}

	*m = *session
	return nil
}
