package models

import (
	"log"
	"messenger/internal/app/errors"
	"messenger/internal/infrastructure/validators"
	ctx "messenger/internal/infrastructure/handler_context"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthReqModel struct {
	Username string
	Password string
}

func (m *AuthReqModel) FromRequest(body ctx.RequestBody) *errors.Error {
	err := new(errors.Error)

	if m.Username, err = body.GetString("username"); err != nil {
		return err
	}
	if m.Password, err = body.GetString("password"); err != nil {
		return err
	}

	if err := validators.String(m.Username).
		MaxLen(50).
		MinLen(3).
		Validate(); err != nil {
		return err.WithFieldName("username").BuildError()
	}

	if err := validators.String(m.Password).
		MaxLen(50).
		MinLen(3).
		Validate(); err != nil {
		return err.WithFieldName("password").BuildError()
	}

	return nil
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

type UserDbModelCollection []*UserDbModel

func (c *UserDbModelCollection) FromDb(rows pgx.Rows) *errors.Error {
	for rows.Next() {
		model := new(UserDbModel)
		if err := model.FromDb(rows); err != nil {
			return err.BuildError()
		}
		*c = append(*c, model)
	}
	return nil
}

func (c UserDbModelCollection) ToResponse() []map[string]any {
	response := make([]map[string]any, len(c))
	for _, model := range c {
		response = append(response, model.ToResponse())
	}
	return response
}
