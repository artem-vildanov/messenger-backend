package repo

import (
	"context"
	"messenger/internal/app/errors"
	"messenger/internal/app/models"
)

type AuthRepository interface {
	GetSession(ctx context.Context, sessionId string) (*models.SessionModel, *errors.Error)
	SaveSession(ctx context.Context, session *models.SessionModel) *errors.Error
	DeleteSession(ctx context.Context, sessionId string) *errors.Error
	DeleteAllExpired(ctx context.Context) *errors.Error
}

type AuthRepositoryImpl struct {
	Repository
}

func (r *AuthRepositoryImpl) GetSession(ctx context.Context, sessionId string) (*models.SessionModel, *errors.Error) {
	sql := `
		select * from sessions
		where id = $1;
	`

	sessionModel := new(models.SessionModel)
	if err := sessionModel.FromDb(r.queryRow(ctx, sql, sessionId)); err != nil {
		return nil, err.WithId(sessionId).BuildError()
	}

	return sessionModel, nil
}

func (u *AuthRepositoryImpl) SaveSession(ctx context.Context, session *models.SessionModel) *errors.Error {
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

func (u *AuthRepositoryImpl) DeleteSession(ctx context.Context, sessionId string) *errors.Error {
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
func (u *AuthRepositoryImpl) DeleteAllExpired(ctx context.Context) *errors.Error {
	return nil
}
