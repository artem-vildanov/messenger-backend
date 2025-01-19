package storage

import (
	"context"
	"errors"
	"messenger/internal/domain/models"
	appErrors "messenger/internal/infrastructure/errors"

	"github.com/jmoiron/sqlx"
)

const (
	userNotFound appErrors.ResponseMessage = "user not found"
)

type UserStorage struct {
	*AbstractStorage[models.UserModel]
}

func NewUserStorage(pgClient *sqlx.DB) *UserStorage {
	return &UserStorage{
		&AbstractStorage[models.UserModel]{pgClient},
	}
}

func (r *UserStorage) GetByUsername(
	ctx context.Context,
	username string,
) (*models.UserModel, error) {
	sql := `
		select * from users
		where username = $1;
	`

	user, err := r.findOne(ctx, sql, username)
	if err != nil {
		if appErrors.WrappedErrorIs(err, appErrors.ErrNotFound) {
			return nil, appErrors.Wrap(
				appErrors.ErrNotFoundWithMessage("user not found"),
				err,
				errors.New("GetByUsername"),
			)
		}
		return nil, appErrors.Wrap(
			err,
			errors.New("GetByUsername"),
		)
	}

	return user, nil
}

func (r *UserStorage) GetById(
	ctx context.Context,
	id int,
) (*models.UserModel, error) {
	sql := `
		select * from users
		where id = $1;
	`

	user, err := r.findOne(ctx, sql, id)
	if err != nil {
		if appErrors.WrappedErrorIs(err, appErrors.ErrNotFound) {
			return nil, appErrors.Wrap(
				appErrors.ErrNotFoundWithMessage("user not found"),
				err,
				errors.New("GetById"),
			)
		}
		return nil, appErrors.Wrap(
			appErrors.ErrInternal,
			err,
			errors.New("GetById"),
		)
	}
	return user, nil
}

func (r *UserStorage) Create(
	ctx context.Context,
	authModel *models.AuthModel,
) (int, error) {
	sql := `
		insert into users (username, password_hash)
		values ($1, $2)
		returning id;
	`

	var userId int
	if err := r.queryRow(
		ctx,
		sql,
		authModel.Username,
		authModel.Password,
	).Scan(&userId); err != nil {
		if appErrors.IsUniqueViolationErr(err) {
			return 0, appErrors.Wrap(
				appErrors.ErrBadRequestWithMessage("user already exists"),
				err,
				errors.New("Create"),
			)
		}
		return 0, appErrors.Wrap(
			appErrors.ErrInternal,
			err,
			errors.New("Create"),
		)
	}

	return userId, nil
}

func (r *UserStorage) GetAll(ctx context.Context) (
	[]*models.UserModel,
	error,
) {
	sql := `
		select * from users;
	`

	users, err := r.findSlice(ctx, sql)
	if err != nil {
		return nil, appErrors.Wrap(
			err,
			errors.New("GetAll"),
		)
	}

	return users, nil
}
