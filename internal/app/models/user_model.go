package models

import (
	"log"
	"messenger/internal/app/errors"
	"messenger/internal/app/handlers/ctx"
	"messenger/internal/app/models/validator"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthReqModel struct {
	Username string
	Password string
}

func NewAuthReqModel(body ctx.RequestBody) (*AuthReqModel, *errors.Error) {
	model := &AuthReqModel{
		body["username"].(string),
		body["password"].(string),
	}

	if err := validator.
		String(model.Username).
		MaxLen(50).
		MinLen(3).
		Validate(); err != nil {
		return nil, err.WithFieldName("username").BuildError()
	}

	if err := validator.
		String(model.Password).
		MaxLen(50).
		MinLen(3).
		Validate(); err != nil {
		return nil, err.WithFieldName("password").BuildError()
	}

	return model, nil
}

func (a *AuthReqModel) HashPassword() *errors.Error {
	hash, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %s", err.Error())
		return errors.InternalError()
	}
	a.Password = string(hash)
	return nil
}

func (a *AuthReqModel) VerifyPassword(passwordHash string) *errors.Error {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(a.Password),
	); err != nil {
		return errors.ForbiddenError().WithVerbose("wrong credentials")
	}
	return nil
}

type UserDbModel struct {
	ID           int
	Username     string
	PasswordHash string
}

func (m *UserDbModel) FromDb(row pgx.Row) *errors.UserError {
	if err := row.Scan(
		&m.ID,
		&m.Username,
		&m.PasswordHash,
	); err != nil {
		if pgx.ErrNoRows.Error() == err.Error() {
			return errors.UserDoesntExistsError()
		}
		return errors.FailedToFindUserError().WithReason(err.Error())
	}
	return nil
}

func (m *UserDbModel) ToResponse() map[string]any {
	return map[string]any{
		"id":       m.ID,
		"username": m.Username,
	}
}
