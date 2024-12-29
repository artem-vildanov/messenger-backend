package repository

import (
	"context"
	"messenger/internal/app/errors"
	"messenger/internal/app/models"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*models.UserDbModel, *errors.Error)
	GetById(ctx context.Context, id int) (*models.UserDbModel, *errors.Error)
	ExistsByUsername(ctx context.Context, username string) *errors.Error
	Create(ctx context.Context, userReqModel *models.AuthReqModel) (int, *errors.Error)
	GetAll(ctx context.Context) (models.UserDbModelCollection, *errors.Error)
}

type UserRepositoryImpl struct {
	AbstractRepository
}

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*models.UserDbModel, *errors.Error) {
	sql := `
		select * from users
		where username = $1;
	`
	userDbModel := new(models.UserDbModel)
	if err := userDbModel.FromDb(
		r.queryRow(ctx, sql, username),
	); err != nil {
		return nil, err.WithName(username).BuildError()
	}
	return userDbModel, nil
}

func (r *UserRepositoryImpl) GetById(ctx context.Context, id int) (*models.UserDbModel, *errors.Error) {
	sql := `
		select * from users
		where id = $1;
	`
	userDbModel := new(models.UserDbModel)
	if err := userDbModel.FromDb(
		r.queryRow(ctx, sql, id),
	); err != nil {
		return nil, err.WithId(id).BuildError()
	}
	return userDbModel, nil
}

func (r *UserRepositoryImpl) ExistsByUsername(ctx context.Context, username string) *errors.Error {
	sql := `
		select exists (
			select 1
			from users
			where username = $1
		);
	`

	var exists bool
	if err := r.queryRow(ctx, sql, username).Scan(&exists); err != nil {
		return errors.FailedToFindUserError().
			WithName(username).
			WithReason(err.Error()).
			BuildError()
	}

	if !exists {
		return errors.UserDoesntExistsError().
			WithName(username).
			BuildError()
	}

	return nil
}

func (r *UserRepositoryImpl) Create(ctx context.Context, userReqModel *models.AuthReqModel) (int, *errors.Error) {
	sql := `
		insert into users (username, password_hash)
		values ($1, $2)
		returning id;
	`

	var userId int
	if err := r.queryRow(
		ctx,
		sql,
		userReqModel.Username,
		userReqModel.Password,
	).Scan(&userId); err != nil {
		return 0, errors.HandleCreateUserError(err).
			WithName(userReqModel.Username).
			BuildError()
	}

	return userId, nil
}

func (r *UserRepositoryImpl) GetAll(ctx context.Context) (models.UserDbModelCollection, *errors.Error) {
	sql := `
		select * from users;
	`

	models := make(models.UserDbModelCollection, 0)
	rows, err := r.query(ctx, sql)
	if err != nil {
		// todo implement
	}
	
	if err := models.FromDb(rows); err != nil {
		return nil, err
	}
	
	return models, nil
}
