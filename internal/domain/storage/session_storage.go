package storage

import (
	"context"
	"messenger/internal/domain/models"
	"messenger/internal/infrastructure/errors"

	"github.com/jmoiron/sqlx"
)

type SessionStorage struct {
	*AbstractStorage[models.SessionModel]
}

func NewSessionStorage(pg *sqlx.DB) *SessionStorage {
	return &SessionStorage{
		&AbstractStorage[models.SessionModel]{pg},
	}
}

func (r *SessionStorage) GetSessionById(
	requestContext context.Context,
	sessionId string,
) (*models.SessionModel, *errors.Error) {
	sql := `
		select * from sessions
		where id = $1;
	`

	session, err := r.findOne(requestContext, sql, sessionId)
	if err != nil {
		return nil, r.handleSessionNotFoundError(err).
			WithField("sessionId", sessionId)
	}

	return session, nil
}

func (r *SessionStorage) GetSessionByUserId(
	ctx context.Context,
	userId int,
) (*models.SessionModel, *errors.Error) {
	sql := `
		select * from sessions
		where user_id = $1;
	`

	session, err := r.findOne(ctx, sql, userId)
	if err != nil {
		return nil, r.handleSessionNotFoundError(err).
			WithField("userId", userId)
	}

	return session, nil
}

func (u *SessionStorage) SaveSession(ctx context.Context, session *models.SessionModel) *errors.Error {
	sql := `
		insert into sessions(id, user_id, expires_at)
		values ($1, $2, $3);
	`
	if err := u.exec(ctx, sql, session.Id, session.UserId, session.ExpiresAt); err != nil {
		return errors.InternalError().
			WithLogMessage(err.Error(), "failed to save session").
			WithField("Session", session)
	}

	return nil
}

func (u *SessionStorage) DeleteSession(ctx context.Context, sessionId string) *errors.Error {
	sql := `
		delete from sessions
		where id = $1;
	`

	if err := u.exec(ctx, sql, sessionId); err != nil {
		return errors.InternalError().
			WithLogMessage(err.Error(), "failed to delete session").
			WithField("sessionId", sessionId)
	}

	return nil
}

// todo
func (u *SessionStorage) DeleteAllExpired(ctx context.Context) *errors.Error {
	return nil
}

func (r *SessionStorage) handleSessionNotFoundError(err *errors.Error) *errors.Error {
	if err.ResponseMessage == errors.NotFoundMessage {
		return errors.UnauthorizedError().
			WithLogMessage("session not found")
	} else {
		return err.WithLogMessage("failed to find session")
	}
}
