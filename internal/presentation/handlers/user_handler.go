package handlers

import (
	"context"
	"messenger/internal/domain/models"
	"messenger/internal/infrastructure/errors"
	ctx "messenger/internal/infrastructure/utils/handler_utils"
	"messenger/internal/presentation/dto"
)

type UserGetter interface {
	GetById(ctx context.Context, id int) (*models.UserModel, *errors.Error)
	GetAll(ctx context.Context) ([]*models.UserModel, *errors.Error)
}

type UserHandler struct {
	userRepo UserGetter
}

func NewUserHandler(userRepo UserGetter) *UserHandler {
	return &UserHandler{userRepo}
}

func (u *UserHandler) GetUserById(handlerContext *ctx.HandlerContext) *errors.Error {
	userId, err := handlerContext.PathParams.GetInteger("userId")
	if err != nil {
		return err
	}

	user, err := u.userRepo.GetById(
		handlerContext.Request.Context(), 
		userId,
	)
	if err != nil {
		return err
	}

	handlerContext.Response().
		WithContent(dto.NewUserResponse(user)).
		Json()

	return nil
}

func (u *UserHandler) GetMyUser(handlerContext *ctx.HandlerContext) *errors.Error {
	user, err := u.userRepo.GetById(
		handlerContext.Request.Context(), 
		handlerContext.AuthUserId,
	)
	if err != nil {
		return err
	}

	handlerContext.Response().
		WithContent(dto.NewUserResponse(user)).
		Json()

	return nil
}

func (u *UserHandler) GetAllUsers(handlerContext *ctx.HandlerContext) *errors.Error {
	models, err := u.userRepo.GetAll(handlerContext.Request.Context())
	if err != nil {
		return err
	}

	handlerContext.Response().
		WithContent(dto.NewMultipleUsersResponse(models)).
		Json()

	return nil
}
