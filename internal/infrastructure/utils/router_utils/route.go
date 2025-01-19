package router_utils

import (
	"messenger/internal/infrastructure/utils/handler_utils"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler = func(*handler_utils.HandlerContext) error

func toHttpHandler(handler Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		handlerContext := &handler_utils.HandlerContext{
			Request:        request,
			ResponseWriter: responseWriter,
			QueryParams:    request.URL.Query(),
			PathParams:     mux.Vars(request),
		}

		if err := handler(handlerContext); err != nil {
			handlerContext.ErrorResponse(err)
		}
	})
}

type Route struct {
	method  Method
	path    string
	handler Handler
}

func NewRoute(method Method, path string, handler Handler) *Route {
	return &Route{
		method,
		path,
		handler,
	}
}

func (r *Route) Middleware(middlewares ...Middleware) *Route {
	for _, middleware := range middlewares {
		next := r.handler
		r.handler = Handler(func(handlerContext *handler_utils.HandlerContext) error {
			return middleware(handlerContext, next)
		})
	}
	return r
}

type Method string

const (
	Get    Method = "GET"
	Post   Method = "POST"
	Put    Method = "PUT"
	Delete Method = "DELETE"
)
