package storage

import (
	"context"
	"errors"
	"messenger/internal/domain/models"
	appErrors "messenger/internal/infrastructure/errors"

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
) (*models.SessionModel, error) {
	sql := `
		select * from sessions
		where id = $1;
	`

	session, err := r.findOne(requestContext, sql, sessionId)
	if err != nil {
		if appErrors.WrappedErrorIs(err, appErrors.ErrNotFound) {
			return nil, appErrors.Wrap(
				appErrors.ErrSessionExpired,
				err,
				errors.New("GetSessionById"),
			)
		}
		return nil, appErrors.Wrap(
			err, 
			errors.New("GetSessionById"),
		)
	}

	return session, nil
}

func (r *SessionStorage) GetSessionByUserId(
	ctx context.Context,
	userId int,
) (*models.SessionModel, error) {
	sql := `
		select * from sessions
		where user_id = $1;
	`

	session, err := r.findOne(ctx, sql, userId)
	if err != nil {
		if appErrors.WrappedErrorIs(err, appErrors.ErrNotFound) {
			return nil, appErrors.Wrap(
				appErrors.ErrUnauthorized,
				err,
				errors.New("GetSessionById"),
			)
		}
		return nil, appErrors.Wrap(
			err, 
			errors.New("GetSessionById"),
		)

	}

	return session, nil
}

func (u *SessionStorage) SaveSession(
	ctx context.Context, 
	session *models.SessionModel,
) error {
	sql := `
		insert into sessions(id, user_id, expires_at)
		values ($1, $2, $3);
	`
	if err := u.exec(
		ctx, 
		sql, 
		session.Id, 
		session.UserId, 
		session.ExpiresAt,
	); err != nil {
		return appErrors.Wrap(
			appErrors.ErrInternal,
			err,
			errors.New("SaveSession"),
		)
	}

	return nil
}

func (u *SessionStorage) DeleteSession(
	ctx context.Context, 
	sessionId string,
) error {
	sql := `
		delete from sessions
		where id = $1;
	`

	if err := u.exec(ctx, sql, sessionId); err != nil {
		return appErrors.Wrap(
			appErrors.ErrInternal,
			err,
			errors.New("DeleteSession"),
		)
	}

	return nil
}

// todo
func (u *SessionStorage) DeleteAllExpired(ctx context.Context) error {
	return nil
}
