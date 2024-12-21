package errors

import (
	"fmt"
)

type UserError struct {
	*Error
	errorReason string
	username    string
	userId      int
}

func (u *UserError) WithName(username string) *UserError {
	u.username = username
	return u
}

func (u *UserError) WithId(id int) *UserError {
	u.userId = id
	return u
}

func (u *UserError) WithReason(reason string) *UserError {
	u.errorReason = reason
	return u
}

func (u *UserError) BuildError() *Error {
	if len(u.username) > 0 {
		u.WithVerbose(fmt.Sprintf("with username: [%s]", u.username))
	}
	if u.userId != 0 {
		u.WithVerbose(fmt.Sprintf("with id: [%d]", u.userId))
	}
	if len(u.errorReason) > 0 {
		u.WithVerbose(fmt.Sprintf("for reason: [%s]", u.errorReason))
	}
	return u.Error
}

func UserAlreadyExistsError() *UserError {
	return &UserError{
		Error: BadRequestError().WithVerbose("user already exists"),
	}
}

func UserDoesntExistsError() *UserError {
	return &UserError{
		Error: BadRequestError().WithVerbose("user doesnt exist"),
	}
}

func FailedToFindUserError() *UserError {
	return &UserError{
		Error: InternalError().WithVerbose("failed to find user"),
	}
}

func FailedToCreateUserError() *UserError {
	return &UserError{
		Error: InternalError().WithVerbose("failed to create user"),
	}
}
