package ctx

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"messenger/internal/app/errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const SessionIdKey string = "sessionId"
const SessionKey string = "session"
type ContextKey string

type RequestBody map[string]interface{}

func (r *RequestBody) fromRequest(request *http.Request) *errors.Error {
	if request.ContentLength == 0 {
		return nil
	}

	defer request.Body.Close()
	rawBody, err := io.ReadAll(request.Body)
	if err != nil {
		return errors.InternalError().WithVerbose(fmt.Sprintf("Error reading body: %s", err))
	}

	err = json.Unmarshal(rawBody, r)
	if err != nil {
		return errors.InternalError().WithVerbose(fmt.Sprintf("Error parsing JSON: %s", err))
	}

	return nil
}

type PathParams map[string]string

func (pathParams PathParams) GetInteger(key string) (int, *errors.Error) {
	strValue := pathParams[key]
	if len(strValue) == 0 {
		return 0, errors.NotFoundError()
	}
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, errors.NotFoundError()
	}
	return intValue, nil
}

type HandlerContext struct {
	Request        *http.Request
	responseWriter http.ResponseWriter
	PathParams     PathParams
	Body           RequestBody
}

func NewHandlerContext(
	responseWriter http.ResponseWriter,
	request *http.Request,
) (*HandlerContext, *errors.Error) {
	var body RequestBody
	if err := body.fromRequest(request); err != nil {
		return nil, err
	}

	return &HandlerContext{
		request,
		responseWriter,
		mux.Vars(request),
		body,
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

func (c *HandlerContext) ErrorResponse(err *errors.Error) {
	c.Response().
		WithCode(err.GetCode()).
		WithContent(err.Error()).
		Json()
}

func (c *HandlerContext) Response() *ResponseBuilder {
	c.responseWriter.Header().Set("Content-Type", "application/json")
	return &ResponseBuilder{
		responseWriter: c.responseWriter,
		code: 200,
		content: "OK",
	}
}
