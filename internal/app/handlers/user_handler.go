package handlers

import (
	"messenger/internal/app/errors"
	repo "messenger/internal/app/repository"
	ctx "messenger/internal/infrastructure/handler_context"
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

	handlerContext.Response().
		WithContent(userDbModel.ToResponse()).
		Json()

	return nil
}

func (u *UserHandler) GetMyUser(handlerContext *ctx.HandlerContext) *errors.Error {
	requestContext := handlerContext.Request.Context()

	userDbModel, err := u.userRepo.GetById(requestContext, handlerContext.Session.UserId)
	if err != nil {
		return err
	}

	handlerContext.Response().
		WithContent(userDbModel.ToResponse()).
		Json()
	
	return nil
}

func (u *UserHandler) GetAllUsers(handlerContext *ctx.HandlerContext) *errors.Error {
	requestContext := handlerContext.Request.Context()

	models, err := u.userRepo.GetAll(requestContext)
	if err != nil {
		return err
	}

	handlerContext.Response().
		WithContent(models.ToResponse()).
		Json()
	
	return nil
}
