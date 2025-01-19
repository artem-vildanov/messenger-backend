package handlers

import (
	"context"
	"errors"
	"messenger/internal/domain/models"
	"messenger/internal/infrastructure/config"
	appErrors "messenger/internal/infrastructure/errors"
	ctx "messenger/internal/infrastructure/utils/handler_utils"
	"messenger/internal/infrastructure/utils/mapping_utils"
	"messenger/internal/presentation/dto"
)

type SessionService interface {
	Authorize(
		ctx context.Context,
		loginModel *models.AuthModel,
	) (*models.SessionModel, error)
	Signup(
		ctx context.Context,
		signupModel *models.AuthModel,
	) (*models.SessionModel, error)
}

type SessionStorage interface {
	DeleteSession(ctx context.Context, sessionId string) error
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

func (h *AuthHandler) Register(handlerContext *ctx.HandlerContext) error {
	signupRequest, err := mapping_utils.FromRequest[*dto.AuthRequest](
		handlerContext.Request,
	)
	if err != nil {
		return appErrors.Wrap(err, errors.New("Register"))
	}

	session, err := h.sessionService.Signup(
		handlerContext.Request.Context(),
		signupRequest.ToDomain(),
	)
	if err != nil {
		return appErrors.Wrap(err, errors.New("Register"))
	}

	handlerContext.Response().
		WithCookie(session.Id).
		WithContent(map[string]any{
			"userId": session.UserId,
		}).
		Json()

	return nil
}

func (h *AuthHandler) Login(handlerContext *ctx.HandlerContext) error {
	loginModel, err := mapping_utils.FromRequest[*dto.AuthRequest](
		handlerContext.Request,
	)
	if err != nil {
		return appErrors.Wrap(err, errors.New("Login"))
	}

	session, err := h.sessionService.Authorize(
		handlerContext.Request.Context(),
		loginModel.ToDomain(),
	)
	if err != nil {
		return appErrors.Wrap(err, errors.New("Login"))
	}

	handlerContext.Response().
		WithCookie(session.Id).
		WithContent(map[string]any{
			"userId": session.UserId,
		}).
		Json()

	return nil
}

func (h *AuthHandler) Logout(handlerContext *ctx.HandlerContext) error {
	if err := h.sessionDeleter.DeleteSession(
		handlerContext.Request.Context(),
		handlerContext.SessionId,
	); err != nil {
		return appErrors.Wrap(err, errors.New("Logout"))
	}

	handlerContext.Response().Empty()

	return nil
}
