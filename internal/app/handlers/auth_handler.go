package handlers

import (
	"messenger/internal/app/errors"
	"messenger/internal/app/models"
	repo "messenger/internal/app/repository"
	"messenger/internal/infrastructure/config"
	ctx "messenger/internal/infrastructure/handler_context"
)

type AuthHandler struct {
	userRepo repo.UserRepository
	authRepo repo.SessionRepository
	env      *config.Env
}

func (h *AuthHandler) Construct(
	userRepo repo.UserRepository,
	authRepo repo.SessionRepository,
	env *config.Env,
) {
	h.userRepo = userRepo
	h.authRepo = authRepo
	h.env = env
}

func (h *AuthHandler) Register(handlerContext *ctx.HandlerContext) *errors.Error {
	requestContext := handlerContext.Request.Context()
	var err *errors.Error

	registerModel := new(models.AuthReqModel)
	if err = registerModel.FromRequest(handlerContext.Body); err != nil {
		return err
	}

	registerModel.HashPassword()
	userId, err := h.userRepo.Create(requestContext, registerModel)
	if err != nil {
		return err
	}

	session := ctx.NewSession(userId, h.env.GetSessionTTL())
	if err := h.authRepo.SaveSession(requestContext, session); err != nil {
		return err
	}

	return handlerContext.Response().
		WithCookie(session.ID).
		TextPlain()
}

func (h *AuthHandler) Login(handlerContext *ctx.HandlerContext) *errors.Error {
	requestContext := handlerContext.Request.Context()
	var err *errors.Error

	loginModel := new(models.AuthReqModel)
	if err = loginModel.FromRequest(handlerContext.Body); err != nil {
		return err
	}

	user, err := h.userRepo.GetByUsername(requestContext, loginModel.Username)
	if err != nil {
		return err
	}

	if err := loginModel.VerifyPassword(user.PasswordHash); err != nil {
		return err
	}

	session := ctx.NewSession(user.ID, h.env.GetSessionTTL())
	if err := h.authRepo.SaveSession(requestContext, session); err != nil {
		return err
	}

	return handlerContext.Response().
		WithCookie(session.ID).
		TextPlain()
}

func (a *AuthHandler) Logout(handlerContext *ctx.HandlerContext) *errors.Error {
	requestContext := handlerContext.Request.Context()

	if err := a.authRepo.DeleteSession(requestContext, handlerContext.Session.ID); err != nil {
		return err
	}

	return handlerContext.Response().Empty()
}
