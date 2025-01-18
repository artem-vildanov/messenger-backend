package storage

import (
	"context"
	"messenger/internal/domain/models"
	"messenger/internal/infrastructure/errors"

	"github.com/jmoiron/sqlx"
)

const (
	userNotFound errors.ResponseMessage = "user not found"
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
) (*models.UserModel, *errors.Error) {
	sql := `
		select * from users
		where username = $1;
	`

	user, err := r.findOne(ctx, sql, username)
	if err != nil {
		return nil, err.WithResponseMessage(userNotFound).WithField("username", username)
	}

	return user, nil
}

func (r *UserStorage) GetById(ctx context.Context, id int) (*models.UserModel, *errors.Error) {
	sql := `
		select * from users
		where id = $1;
	`

	user, err := r.findOne(ctx, sql, id)
	if err != nil {
		return nil, err.WithResponseMessage(userNotFound).WithField("id", id)
	}
	return user, nil
}

func (r *UserStorage) Create(
	ctx context.Context, 
	authModel *models.AuthModel,
) (int, *errors.Error) {
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
		if r.isUniqueViolation(err) {
			return 0, errors.BadRequestError().
				WithLogMessage(err.Error()).
				WithResponseMessage("user already exists").
				WithField("AuthModel", authModel)
		}
		return 0, errors.InternalError().
			WithLogMessage(err.Error(), "failed to create user").
			WithField("AuthRequest", authModel)
	}

	return userId, nil
}

func (r *UserStorage) GetAll(ctx context.Context) ([]*models.UserModel, *errors.Error) {
	sql := `
		select * from users;
	`

	users, err := r.findSlice(ctx, sql)
	if err != nil {
		return nil, err.WithLogMessage("failed to get all users")
	}

	return users, nil
}
