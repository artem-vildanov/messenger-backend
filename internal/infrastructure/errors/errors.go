package errors

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/lib/pq"
)

func Wrap(err error, wrappings ...error) error {
	for _, wrapping := range wrappings {
		err = fmt.Errorf("%s: %w", wrapping, err)
	}
	return err
}

func Unwrap(err error) (*UnwrappedError, bool) {
	wrappings := make([]error, 0)
	unwrapped := unwrapRecursive(err, &wrappings)
	if unwrapped, ok := unwrapped.(Error); ok {
		return &UnwrappedError{
			unwrapped,
			wrappings,
		}, true
	}
	return nil, false
}

func unwrapRecursive(toUnwrap error, wrappings *[]error) error {
	wrapping := errors.Unwrap(toUnwrap)
	if wrapping != nil {
		*wrappings = append(*wrappings, wrapping)
		return unwrapRecursive(wrapping, wrappings)
	}
	return toUnwrap
}

type UnwrappedError struct {
	Error
	Wrappings []error
}

func (e UnwrappedError) LogStdout() {
	temp := struct {
		Err       Error    `json:"error"`
		Wrappings []string `json:"wrappings"`
	}{
		Err: e.Error,
	}

	for _, wrapping := range e.Wrappings {
		temp.Wrappings = append(temp.Wrappings, wrapping.Error())
	}

	jsonData, err := json.MarshalIndent(temp, "", "  ")
	if err != nil {
		log.Println("Failed to cast Error into JSON", err)
		return
	}

	log.Println(string(jsonData))
}

type ResponseMessage string

func NewError(code int, reason ResponseMessage) Error {
	return Error{
		Code:            code,
		ResponseMessage: reason,
	}
}

type Error struct {
	Code            int             `json:"code"`
	ResponseMessage ResponseMessage `json:"responseMessage"`
}

func (e Error) Error() string {
	return string(e.ResponseMessage)
}

var (
	ErrBadRequest     = NewError(http.StatusBadRequest, "bad request")
	ErrNotFound       = NewError(http.StatusNotFound, "not found")
	ErrUnauthorized   = NewError(http.StatusUnauthorized, "unauthorized")
	ErrSessionExpired = NewError(http.StatusUnauthorized, "session expired")
	ErrForbidden      = NewError(http.StatusForbidden, "forbidden")
	ErrInternal       = NewError(http.StatusInternalServerError, "internal error")
)

func ErrBadRequestWithMessage(message string) error {
	return NewError(
		http.StatusBadRequest,
		ResponseMessage(message),
	)
}

func ErrNotFoundWithMessage(message string) error {
	return NewError(
		http.StatusNotFound,
		ResponseMessage(message),
	)
}

func IsUniqueViolationErr(err error) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func IsNoRowsErr(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func WrappedErrorIs(wrapped error, target error) bool {
	unwrapped, ok := Unwrap(wrapped)
	if !ok {
		return false
	}

	return errors.Is(unwrapped.Error, target)
}
