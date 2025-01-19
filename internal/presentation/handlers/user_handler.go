package handlers

import (
	"context"
	"errors"
	"messenger/internal/domain/models"
	appErrors "messenger/internal/infrastructure/errors"
	ctx "messenger/internal/infrastructure/utils/handler_utils"
	"messenger/internal/presentation/dto"
)

type UserGetter interface {
	GetById(ctx context.Context, id int) (*models.UserModel, error)
	GetAll(ctx context.Context) ([]*models.UserModel, error)
}

type UserHandler struct {
	userRepo UserGetter
}

func NewUserHandler(userRepo UserGetter) *UserHandler {
	return &UserHandler{userRepo}
}

func (u *UserHandler) GetUserById(handlerContext *ctx.HandlerContext) error {
	userId, err := handlerContext.PathParams.GetInteger("userId")
	if err != nil {
		return appErrors.Wrap(err, errors.New("GetUserById"))
	}

	user, err := u.userRepo.GetById(
		handlerContext.Request.Context(), 
		userId,
	)
	if err != nil {
		return appErrors.Wrap(err, errors.New("GetUserById"))
	}

	handlerContext.Response().
		WithContent(dto.NewUserResponse(user)).
		Json()

	return nil
}

func (u *UserHandler) GetMyUser(handlerContext *ctx.HandlerContext) error {
	user, err := u.userRepo.GetById(
		handlerContext.Request.Context(), 
		handlerContext.AuthUserId,
	)
	if err != nil {
		return appErrors.Wrap(err, errors.New("GetMyUser"))
	}

	handlerContext.Response().
		WithContent(dto.NewUserResponse(user)).
		Json()

	return nil
}

func (u *UserHandler) GetAllUsers(handlerContext *ctx.HandlerContext) error {
	models, err := u.userRepo.GetAll(handlerContext.Request.Context())
	if err != nil {
		return appErrors.Wrap(err, errors.New("GetAllUsers"))
	}

	handlerContext.Response().
		WithContent(dto.NewMultipleUsersResponse(models)).
		Json()

	return nil
}
