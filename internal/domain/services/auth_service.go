package services

import (
	"context"
	"errors"
	"messenger/internal/domain/models"
	"messenger/internal/infrastructure/config"
	appErrors "messenger/internal/infrastructure/errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SessionStorage interface {
	GetSessionById(
		requestContext context.Context,
		sessionId string,
	) (*models.SessionModel, error)
	DeleteSession(ctx context.Context, sessionId string) error
	SaveSession(ctx context.Context, session *models.SessionModel) error
}

type UserStorage interface {
	Create(ctx context.Context, userReqModel *models.AuthModel) (int, error)
	GetByUsername(
		ctx context.Context,
		username string,
	) (*models.UserModel, error)
}

type SessionService struct {
	sessionStorage SessionStorage
	userStorage    UserStorage
	env            *config.Env
}

func NewSessionService(
	sessionStorage SessionStorage,
	userCreator UserStorage,
	env *config.Env,
) *SessionService {
	return &SessionService{
		sessionStorage,
		userCreator,
		env,
	}
}

func (s *SessionService) AuthenticateBySessionId(
	ctx context.Context,
	sessionId string,
) (*models.SessionModel, error) {
	session, err := s.sessionStorage.GetSessionById(ctx, sessionId)
	if err != nil {
		return nil, appErrors.Wrap(
			err,
			errors.New("AuthenticateBySessionId"),
		)
	}

	if session.ExpiresAt.Before(time.Now()) {
		if err := s.sessionStorage.DeleteSession(ctx, session.Id); err != nil {
			return nil, appErrors.Wrap(
				err,
				errors.New("AuthenticateBySessionId"),
			)
		}

		return nil, appErrors.Wrap(
			appErrors.ErrSessionExpired,
			errors.New("AuthenticateBySessionId"),
		)
	}

	return session, nil
}

func (s *SessionService) Signup(
	ctx context.Context,
	signupModel *models.AuthModel,
) (*models.SessionModel, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(signupModel.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, appErrors.Wrap(
			appErrors.ErrInternal,
			err,
			errors.New("Signup"),
		)
	}
	signupModel.Password = string(hash)

	userId, err := s.userStorage.Create(ctx, signupModel)
	if err != nil {
		return nil, appErrors.Wrap(err, errors.New("Signup"))
	}

	session := s.createSession(userId)
	if err := s.sessionStorage.SaveSession(ctx, session); err != nil {
		return nil, appErrors.Wrap(err, errors.New("Signup"))
	}

	return session, nil
}

func (s *SessionService) Authorize(
	ctx context.Context,
	loginModel *models.AuthModel,
) (*models.SessionModel, error) {
	user, err := s.userStorage.GetByUsername(ctx, loginModel.Username)
	if err != nil {
		return nil, appErrors.Wrap(err, errors.New("Authorize"))
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(loginModel.Password),
	); err != nil {
		return nil, appErrors.Wrap(
			appErrors.ErrBadRequestWithMessage("wrong password"),
			err,
			errors.New("Authorize"),
		)
	}

	session := s.createSession(user.Id)
	if err := s.sessionStorage.SaveSession(ctx, session); err != nil {
		return nil, appErrors.Wrap(err, errors.New("Authorize"))
	}

	return session, nil
}

func (s *SessionService) createSession(userId int) *models.SessionModel {
	return &models.SessionModel{
		Id:     uuid.NewString(),
		UserId: userId,
		ExpiresAt: time.Now().Add(
			time.Duration(s.env.SessionTTL) * time.Minute,
		),
	}
}
