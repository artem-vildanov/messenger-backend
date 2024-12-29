package router

import (
	"log"
	"messenger/internal/app/errors"
	ctx "messenger/internal/infrastructure/handler_context"
	"net/http"
)

type HandlerFunction func(*ctx.HandlerContext) *errors.Error

func (handler HandlerFunction) ToHttpHandler() http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request)  {
		context, err := ctx.NewHandlerContext(responseWriter, request)
		if err != nil {
			context.ErrorResponse(err)
		}

		if err := handler(context); err != nil {
			log.Println(err.GetVerbose())
			context.ErrorResponse(err)
		}
	})
}

type route struct {
	method  Method
	path    string
	handler HandlerFunction
}

func Route(method Method, path string, handler HandlerFunction) *route {
	return &route{
		method,
		path,
		handler,
	}
}

func (r *route) Middleware(middlewares ...middleware) *route {
	for _, middleware := range middlewares {
		next := r.handler
		r.handler = HandlerFunction(func(handlerContext *ctx.HandlerContext) *errors.Error {
			if err := middleware.MiddlewareFunc(handlerContext); err != nil {
				return err
			}
			return next(handlerContext)
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
