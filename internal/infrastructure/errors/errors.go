package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseMessage string

const (
	BadRequestMessage    ResponseMessage = "bad request"
	NotFoundMessage      ResponseMessage = "not found"
	InternalErrorMessage ResponseMessage = "internal error"
	ForbiddenMessage     ResponseMessage = "forbidden"
	UnauthorizedMessage  ResponseMessage = "unauthorized"
)

func NewError(code int, reason ResponseMessage) *Error {
	return &Error{
		Code:            code,
		ResponseMessage: reason,
		LogMessages:     make([]string, 0),
		LogData:         make(map[string]any),
	}
}

type Error struct {
	Code            int             `json:"code"`
	ResponseMessage ResponseMessage `json:"responseMessage"`
	LogMessages     []string        `json:"logMessages"`
	LogData         map[string]any  `json:"logData"`
	OriginalError   error           `json:"originalError"`
}

func (e *Error) Error() string {
	return string(e.ResponseMessage)
}

func (e *Error) WithField(key string, value any) *Error {
	e.LogData[key] = value
	return e
}

func (e *Error) WithOriginalError(err error) *Error {
	e.OriginalError = err
	return e
}

func (e *Error) WithLogMessage(messages ...string) *Error {
	e.LogMessages = append(e.LogMessages, messages...)
	return e
}

func (e *Error) WithResponseMessage(message ResponseMessage) *Error {
	e.ResponseMessage = message
	e.WithLogMessage(string(message))
	return e
}

func (e *Error) LogStdout() {
	jsonData, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		log.Println("Failed to cast Error into json", err)
		return
	}

	log.Println(string(jsonData))
}

func BadRequestError() *Error {
	return NewError(http.StatusBadRequest, BadRequestMessage)
}

func NotFoundError() *Error {
	return NewError(http.StatusNotFound, NotFoundMessage)
}

func UnauthorizedError() *Error {
	return NewError(http.StatusUnauthorized, UnauthorizedMessage)
}

func ForbiddenError() *Error {
	return NewError(http.StatusForbidden, ForbiddenMessage)
}

func InternalError() *Error {
	return NewError(http.StatusInternalServerError, InternalErrorMessage)
}
