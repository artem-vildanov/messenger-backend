package router

import (
	"log"
	"messenger/internal/app/errors"
	"messenger/internal/app/handlers/ctx"
	"messenger/internal/app/middlewares"
	"net/http"
)

type HandlerFunction func(*ctx.HandlerContext) *errors.Error

/**
func (handler HandlerFunction) Handle(handlerContext *ctx.HandlerContext) *errors.Error {
	return handler(handlerContext)
}

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
*/

func (handler HandlerFunction) HandleRequest(
	responseWriter http.ResponseWriter,
	request *http.Request,
) *errors.Error {
	context, err := ctx.NewHandlerContext(responseWriter, request)
	if err != nil {
		context.ErrorResponse(err)
		return err
	}

	if err := handler(context); err != nil {
		log.Println(err.GetVerbose())
		context.ErrorResponse(err)
		return err
	}

	return nil
}

type route struct {
	method  Method
	path    string
	handler http.Handler
}

func Route(method Method, path string, handler HandlerFunction) *route {
	return &route{
		method,
		path,
		http.HandlerFunc(
			func(responseWriter http.ResponseWriter, request *http.Request) {
				handler.HandleRequest(responseWriter, request)
			},
		),
	}
}

func (r *route) Middleware(middlewares ...middlewares.Middleware) *route {
	for _, middleware := range middlewares {
		next := r.handler

		r.handler = http.HandlerFunc(
			func(responseWriter http.ResponseWriter, request *http.Request) {
				if err := HandlerFunction(middleware.MiddlewareFunc).
					HandleRequest(responseWriter, request); err != nil {
					return
				}
				next.ServeHTTP(responseWriter, request)
			},
		)
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
