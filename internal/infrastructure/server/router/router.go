package router

import (
	"messenger/internal/app/middlewares"

	"github.com/gorilla/mux"
)

type RoutesGroup struct {
	prefix string
	routes []*route
	routesGroups []*RoutesGroup
	middlewares []middlewares.Middleware
}

func RootGroup() *RoutesGroup {
	return NewGroup("")
}

func NewGroup(prefix string) *RoutesGroup {
	return &RoutesGroup{
		prefix,
		make([]*route, 0),
		make([]*RoutesGroup, 0),
		make([]middlewares.Middleware, 0),
	}
}

func (b *RoutesGroup) WithRoutes(routes... *route) *RoutesGroup {
	b.routes = append(b.routes, routes...)
	return b
}

func (b *RoutesGroup) WithGroups(builder... *RoutesGroup) *RoutesGroup {
	b.routesGroups = append(b.routesGroups, builder...)
	return b
}

func (b *RoutesGroup) WithMiddlewares(middlewares... middlewares.Middleware) *RoutesGroup {
	b.middlewares = append(b.middlewares, middlewares...)
	return b
}

func (b *RoutesGroup) BuildRouter(headRouter *mux.Router) {
	subrouter := headRouter.PathPrefix(b.prefix).Subrouter()
	for _, route := range b.routes {
		subrouter.Handle(
			route.path, 
			route.Middleware(b.middlewares...).handler,
		).Methods(string(route.method))
	}

	if len(b.routesGroups) == 0 {
		return
	}

	for _, group := range b.routesGroups {
		group.middlewares = append(group.middlewares, b.middlewares...)
		group.BuildRouter(subrouter)
	}
}
