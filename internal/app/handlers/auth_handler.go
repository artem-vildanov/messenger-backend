package handlers

import (
	"messenger/internal/app/errors"
	"messenger/internal/app/handlers/ctx"
	"messenger/internal/app/models"
	repo "messenger/internal/app/repository"
	"messenger/internal/infrastructure/config"
)

type AuthHandler struct {
	userRepo repo.UserRepository
	authRepo repo.AuthRepository
	env      *config.Env
}

func (h *AuthHandler) Construct(
	userRepo repo.UserRepository,
	authRepo repo.AuthRepository,
	env *config.Env,
) {
	h.userRepo = userRepo
	h.authRepo = authRepo
	h.env = env
}

func (h *AuthHandler) Register(handlerContext *ctx.HandlerContext) *errors.Error {
	requestContext := handlerContext.Request.Context()
	userReqModel, err := models.NewAuthReqModel(handlerContext.Body)
	if err != nil {
		return err
	}

	err = h.userRepo.ExistsByUsername(requestContext, userReqModel.Username)
	if err != nil {
		return err
	}

	userReqModel.HashPassword()
	userId, err := h.userRepo.Create(requestContext, userReqModel)
	if err != nil {
		return err
	}

	session := models.NewSession(userId, h.env.GetSessionTTL())
	if err := h.authRepo.SaveSession(requestContext, session); err != nil {
		return err
	}

	return handlerContext.Response().
		WithCookie(session.ID).
		TextPlain()
}

func (h *AuthHandler) Login(handlerContext *ctx.HandlerContext) *errors.Error {
	requestContext := handlerContext.Request.Context()
	userReqModel, err := models.NewAuthReqModel(handlerContext.Body)
	if err != nil {
		return err
	}

	user, err := h.userRepo.GetByUsername(requestContext, userReqModel.Username)
	if err != nil {
		return err
	}

	if err := userReqModel.VerifyPassword(user.PasswordHash); err != nil {
		return err
	}

	session := models.NewSession(user.ID, h.env.GetSessionTTL())
	if err := h.authRepo.SaveSession(requestContext, session); err != nil {
		return err
	}

	return handlerContext.Response().
		WithCookie(session.ID).
		TextPlain()
}

func (a *AuthHandler) Logout(handlerContext *ctx.HandlerContext) *errors.Error {
	requestContext := handlerContext.Request.Context()
	session := new(models.SessionModel)
	if err := session.FromContext(requestContext); err != nil {
		return err
	}

	if err := a.authRepo.DeleteSession(requestContext, session.ID); err != nil {
		return err
	}

	return handlerContext.Response().Empty()
}
