package handler_context

import (
	"log"
	"messenger/internal/app/errors"
	"net/http"

	"github.com/gorilla/mux"
)

type ContextKey string

type HandlerContext struct {
	Request        *http.Request
	responseWriter http.ResponseWriter
	PathParams     PathParams
	Body           RequestBody
	Session        *Session
}

func NewHandlerContext(
	responseWriter http.ResponseWriter,
	request *http.Request,
) (*HandlerContext, *errors.Error) {
	var body RequestBody
	if err := body.fromRequest(request); err != nil {
		return nil, err
	}

	session := new(Session)
	_ = session.FromContext(request.Context())

	return &HandlerContext{
		Request:        request,
		responseWriter: responseWriter,
		PathParams:     mux.Vars(request),
		Body:           body,
		Session:        session,
	}, nil
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
	c.responseWriter.Header().Set("Content-Type", "application/json")
	return &ResponseBuilder{
		responseWriter: c.responseWriter,
		code:           http.StatusOK,
		content:        "OK",
	}
}

func (c *HandlerContext) ErrorResponse(err *errors.Error) {
	c.Response().
		WithCode(err.GetCode()).
		WithContent(err.Error()).
		Json()
}
