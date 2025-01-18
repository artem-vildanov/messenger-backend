package services

import (
	"context"
	"messenger/internal/domain/models"
	"messenger/internal/infrastructure/config"
	"messenger/internal/infrastructure/errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SessionStorage interface {
	GetSessionById(
		requestContext context.Context,
		sessionId string,
	) (*models.SessionModel, *errors.Error)
	DeleteSession(ctx context.Context, sessionId string) *errors.Error
	SaveSession(ctx context.Context, session *models.SessionModel) *errors.Error
}

type UserStorage interface {
	Create(ctx context.Context, userReqModel *models.AuthModel) (int, *errors.Error)
	GetByUsername(ctx context.Context, username string) (*models.UserModel, *errors.Error)
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
) (*models.SessionModel, *errors.Error) {
	session, err := s.sessionStorage.GetSessionById(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		if err := s.sessionStorage.DeleteSession(ctx, session.Id); err != nil {
			return nil, err.WithField("Session", session)
		}

		return nil, errors.UnauthorizedError().
			WithResponseMessage("session expired").
			WithField("Session", session)
	}

	return session, nil
}

// returns sessionId
func (s *SessionService) Signup(
	ctx context.Context,
	signupModel *models.AuthModel,
) (*models.SessionModel, *errors.Error) {
	// hash password
	hash, hashErr := bcrypt.GenerateFromPassword([]byte(signupModel.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return nil, errors.InternalError().
			WithLogMessage(hashErr.Error(), "failed to hash password").
			WithField("AuthRequest", signupModel)
	}
	signupModel.Password = string(hash)

	userId, err := s.userStorage.Create(ctx, signupModel)
	if err != nil {
		return nil, err.WithField("AuthRequest", signupModel)
	}

	session := s.createSession(userId)
	if err := s.sessionStorage.SaveSession(ctx, session); err != nil {
		return nil, err.WithField("AuthRequest", signupModel)
	}

	return session, nil
}

// returns sessionId
func (s *SessionService) Authorize(
	ctx context.Context,
	loginModel *models.AuthModel,
) (*models.SessionModel, *errors.Error) {
	user, err := s.userStorage.GetByUsername(ctx, loginModel.Username)
	if err != nil {
		return nil, err.WithField("AuthRequest", loginModel)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(loginModel.Password),
	); err != nil {
		return nil, errors.UnauthorizedError().
			WithResponseMessage("wrong password").
			WithField("AuthRequest", loginModel).
			WithLogMessage(err.Error())
	}

	session := s.createSession(user.Id)
	if err := s.sessionStorage.SaveSession(ctx, session); err != nil {
		return nil, err.WithField("AuthRequest", loginModel)
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
