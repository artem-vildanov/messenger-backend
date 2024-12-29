package repository

import (
	"context"
	"messenger/internal/app/errors"
	ctx "messenger/internal/infrastructure/handler_context"
)

type SessionRepository interface {
	GetSession(ctx context.Context, sessionId string) (*ctx.Session, *errors.Error)
	SaveSession(ctx context.Context, session *ctx.Session) *errors.Error
	DeleteSession(ctx context.Context, sessionId string) *errors.Error
	DeleteAllExpired(ctx context.Context) *errors.Error
}

type SessionRepositoryImpl struct {
	AbstractRepository
}

func (r *SessionRepositoryImpl) GetSession(requestContext context.Context, sessionId string) (*ctx.Session, *errors.Error) {
	sql := `
		select * from sessions
		where id = $1;
	`

	sessionModel := new(ctx.Session)
	if err := sessionModel.FromDb(r.queryRow(requestContext, sql, sessionId)); err != nil {
		return nil, err.WithId(sessionId).BuildError()
	}

	return sessionModel, nil
}

func (u *SessionRepositoryImpl) SaveSession(ctx context.Context, session *ctx.Session) *errors.Error {
	sql := `
		insert into sessions(id, user_id, expires_at)
		values ($1, $2, $3);
	`
	if err := u.exec(ctx, sql, session.ID, session.UserId, session.ExpiredAt); err != nil {
		return errors.FailedToCreateSession().
			WithId(session.ID).
			WithUserId(session.UserId).
			WithReason(err.Error()).
			BuildError()
	}

	return nil
}

func (u *SessionRepositoryImpl) DeleteSession(ctx context.Context, sessionId string) *errors.Error {
	sql := `
		delete from sessions
		where id = $1;
	`

	if err := u.exec(ctx, sql, sessionId); err != nil {
		return errors.FailedToDeleteSession().
			WithId(sessionId).
			WithReason(err.Error()).
			BuildError()
	}

	return nil
}

// todo
func (u *SessionRepositoryImpl) DeleteAllExpired(ctx context.Context) *errors.Error {
	return nil
}
