package handlers

import (
	"messenger/internal/app/errors"
	"messenger/internal/app/handlers/ctx"
	"messenger/internal/app/models"
	repo "messenger/internal/app/repository"
)

type UserHandler struct {
	userRepo repo.UserRepository
}

func (u *UserHandler) Construct(userRepo repo.UserRepository) {
	u.userRepo = userRepo
}

func (u *UserHandler) GetUserById(handlerContext *ctx.HandlerContext) *errors.Error {
	requestContext := handlerContext.Request.Context()
	userId, err := handlerContext.PathParams.GetInteger("userId")
	if err != nil {
		return err
	}

	userDbModel, err := u.userRepo.GetById(requestContext, userId)
	if err != nil {
		return err
	}

	return handlerContext.Response().
		WithContent(userDbModel.ToResponse()).
		Json()
}

func (u *UserHandler) GetMyUser(handlerContext *ctx.HandlerContext) *errors.Error {
	requestContext := handlerContext.Request.Context()
	session := new(models.SessionModel)
	if err := session.FromContext(requestContext); err != nil {
		return err
	}

	userDbModel, err := u.userRepo.GetById(requestContext, session.UserId)
	if err != nil {
		return err
	}

	return handlerContext.Response().
		WithContent(userDbModel.ToResponse()).
		Json()
}
