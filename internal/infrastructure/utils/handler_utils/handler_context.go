package handler_utils

import (
	"log"
	"messenger/internal/infrastructure/errors"
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
func (c *HandlerContext) GetLimitOffset() (int, int, *errors.Error) {
	// get and validate limit
	limitStr := c.QueryParams.Get("limit")
	if len(limitStr) == 0 {
		return 0, 0, errors.BadRequestError().
			WithResponseMessage("limit query param not provided")
	}

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, errors.BadRequestError().
			WithResponseMessage(failedToCastToInt(limitStr))
	}

	if limitInt < minimalLimit {
		return 0, 0, errors.BadRequestError().
			WithResponseMessage(unexpectedLimitValue)
	}

	// get and validate offset
	offsetStr := c.QueryParams.Get("offset")
	if len(offsetStr) == 0 {
		return 0, 0, errors.BadRequestError().
			WithResponseMessage(offsetNotProvided)
	}

	offsetInt, err := strconv.Atoi(offsetStr)
	if err != nil {
		return 0, 0, errors.BadRequestError().
			WithResponseMessage(failedToCastToInt(offsetStr))
	}

	if offsetInt < minimalOffset {
		return 0, 0, errors.BadRequestError().
			WithResponseMessage(unexpectedOffsetValue)
	}

	return limitInt, offsetInt, nil
}

func (c *HandlerContext) SessionCookie() (*http.Cookie, *errors.Error) {
	cookie, err := c.Request.Cookie(SessionIdKey)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, errors.UnauthorizedError()
		} else {
			log.Printf("Failed to get cookie: %v\n", err)
			return nil, errors.InternalError()
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

func (c *HandlerContext) ErrorResponse(err *errors.Error) {
	c.Response().
		WithCode(err.Code).
		WithContent(err.Error()).
		Json()
}
