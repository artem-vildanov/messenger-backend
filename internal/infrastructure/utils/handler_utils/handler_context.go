package handler_utils

import (
	"errors"
	"log"
	appErrors "messenger/internal/infrastructure/errors"
	"net/http"
	"net/url"
	"strconv"
)

const (
	offsetNotProvided     = "offset query param not provided"
	limitNotProvided      = "limit query param not provided"
	unexpectedLimitValue  = "unexpected limit value"
	unexpectedOffsetValue = "unexpected offset value"

	minimalLimit  = 5
	minimalOffset = 0
)

const (
	SessionIdKey string = "sessionId"
	SessionKey   string = "session"
)

type ContextKey string

type HandlerContext struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	PathParams     PathParams
	QueryParams    url.Values
	AuthUserId     int
	SessionId      string
}

// returns limit, offset, error
func (c *HandlerContext) GetLimitOffset() (int, int, error) {
	// get and validate limit
	limitStr := c.QueryParams.Get("limit")
	if len(limitStr) == 0 {
		return 0, 0, appErrors.Wrap(
			appErrors.ErrBadRequestWithMessage("limit request param not provided"),
			errors.New("GetLimitOffset"),
		)
	}

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, appErrors.Wrap(
			appErrors.ErrBadRequestWithMessage("invalid limit"),
			errors.New("GetLimitOffset"),
		)
	}

	if limitInt < minimalLimit {
		return 0, 0, appErrors.Wrap(
			appErrors.ErrBadRequestWithMessage("invalid limit"),
			errors.New("GetLimitOffset"),
		)
	}

	// get and validate offset
	offsetStr := c.QueryParams.Get("offset")
	if len(offsetStr) == 0 {
		return 0, 0, appErrors.Wrap(
			appErrors.ErrBadRequestWithMessage("offset request param not provided"),
			errors.New("GetLimitOffset"),
		)
	}

	offsetInt, err := strconv.Atoi(offsetStr)
	if err != nil {
		return 0, 0, appErrors.Wrap(
			appErrors.ErrBadRequestWithMessage("invalid offset"),
			errors.New("GetLimitOffset"),
		)
	}

	if offsetInt < minimalOffset {
		return 0, 0, appErrors.Wrap(
			appErrors.ErrBadRequestWithMessage("invalid offset"),
			errors.New("GetLimitOffset"),
		)
	}

	return limitInt, offsetInt, nil
}

func (c *HandlerContext) SessionCookie() (*http.Cookie, error) {
	cookie, err := c.Request.Cookie(SessionIdKey)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, appErrors.Wrap(
				appErrors.ErrUnauthorized,
				errors.New("GetSessionCookie"),
			)
		} else {
			return nil, appErrors.Wrap(
				appErrors.ErrInternal,
				errors.New("GetSessionCookie"),
			)
		}
	}
	return cookie, nil
}

func (c *HandlerContext) Response() *ResponseBuilder {
	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	return &ResponseBuilder{
		responseWriter: c.ResponseWriter,
		code:           http.StatusOK,
		content:        "OK",
	}
}

func (c *HandlerContext) ErrorResponse(err error) {
	unwrapped, ok := appErrors.Unwrap(err)
	if !ok {
		log.Printf("failed to unwrap error: %s", err.Error())
	}
	unwrapped.LogStdout()
	c.Response().
		WithCode(unwrapped.Code).
		WithContent(unwrapped.ResponseMessage).
		Json()
}
