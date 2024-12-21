package errors

import (
	"fmt"
	"net/http"
)

const (
	badRequest    errorMessage = "bad request"
	notFound      errorMessage = "not found"
	internalError errorMessage = "internal error"
	forbidden     errorMessage = "forbidden"
	unauthorized  errorMessage = "unauthorized"
)

type errorMessage string

type Error struct {
	code    int
	message errorMessage
	verbose string
}

func (a *Error) GetCode() int {
	return a.code
}

func (a *Error) Error() string {
	return string(a.message)
}

func (e *Error) GetVerbose() string {
	return e.verbose
}

func (e *Error) WithVerbose(verbose string) *Error {
	e.verbose = fmt.Sprintf("%s; %s", e.verbose, verbose)
	return e
}

func BadRequestError() *Error {
	return &Error{
		http.StatusBadRequest,
		badRequest,
		string(badRequest),
	}
}

func NotFoundError() *Error {
	return &Error{
		http.StatusNotFound,
		notFound,
		string(notFound),
	}
}

func UnauthorizedError() *Error {
	return &Error{
		http.StatusUnauthorized,
		unauthorized,
		string(unauthorized),
	}
}

func ForbiddenError() *Error {
	return &Error{
		http.StatusForbidden,
		forbidden,
		string(forbidden),
	}
}

func InternalError() *Error {
	return &Error{
		http.StatusInternalServerError,
		internalError,
		string(internalError),
	}
}
