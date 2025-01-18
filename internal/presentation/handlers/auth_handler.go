package handlers

import (
	"context"
	"messenger/internal/domain/models"
	"messenger/internal/infrastructure/config"
	"messenger/internal/infrastructure/errors"
	ctx "messenger/internal/infrastructure/utils/handler_utils"
	"messenger/internal/infrastructure/utils/mapping_utils"
	"messenger/internal/presentation/dto"
)

const (
	failedToReg    = "failed to register"
	failedToLogin  = "failed to login"
	failedToLogout = "failed to logout"
)

type SessionService interface {
	Authorize(
		ctx context.Context,
		loginModel *models.AuthModel,
	) (*models.SessionModel, *errors.Error)
	Signup(
		ctx context.Context,
		signupModel *models.AuthModel,
	) (*models.SessionModel, *errors.Error)
}

type SessionStorage interface {
	DeleteSession(ctx context.Context, sessionId string) *errors.Error
}

type AuthHandler struct {
	sessionDeleter SessionStorage
	sessionService SessionService
}

func NewAuthHandler(
	sessionDeleter SessionStorage,
	sessionService SessionService,
	env *config.Env,
) *AuthHandler {
	return &AuthHandler{
		sessionDeleter,
		sessionService,
	}
}

func (h *AuthHandler) Register(handlerContext *ctx.HandlerContext) *errors.Error {
	signupRequest, err := mapping_utils.FromRequest[*dto.AuthRequest](handlerContext.Request)
	if err != nil {
		return err.WithLogMessage(failedToReg)
	}

	session, err := h.sessionService.Signup(
		handlerContext.Request.Context(),
		signupRequest.ToDomain(),
	)
	if err != nil {
		return err.WithLogMessage(failedToReg)
	}

	handlerContext.Response().
		WithCookie(session.Id).
		WithContent(map[string]any{
			"userId": session.UserId,
		}).
		Json()

	return nil
}

func (h *AuthHandler) Login(handlerContext *ctx.HandlerContext) *errors.Error {
	loginModel, err := mapping_utils.FromRequest[*dto.AuthRequest](handlerContext.Request)
	if err != nil {
		return err.WithLogMessage(failedToLogin)
	}

	session, err := h.sessionService.Authorize(
		handlerContext.Request.Context(),
		loginModel.ToDomain(),
	)
	if err != nil {
		return err.WithLogMessage(failedToLogin)
	}

	handlerContext.Response().
		WithCookie(session.Id).
		WithContent(map[string]any{
			"userId": session.UserId,
		}).
		Json()

	return nil
}






func (a *AuthHandler) Logout(handlerContext *ctx.HandlerContext) *errors.Error {
	if err := a.sessionDeleter.DeleteSession(
		handlerContext.Request.Context(),
		handlerContext.SessionId,
	); err != nil {
		return err.WithLogMessage(failedToLogout)
	}

	handlerContext.Response().Empty()

	return nil
}
