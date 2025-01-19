package router_utils

import (
	"messenger/internal/infrastructure/utils/handler_utils"

	"github.com/gorilla/mux"
)

type Middleware func(*handler_utils.HandlerContext, Handler) error

type RoutesGroup struct {
	prefix       string
	routes       []*Route
	routesGroups []*RoutesGroup
	middlewares  []Middleware
}

func RootGroup() *RoutesGroup {
	return NewGroup("")
}

func NewGroup(prefix string) *RoutesGroup {
	return &RoutesGroup{
		prefix,
		make([]*Route, 0),
		make([]*RoutesGroup, 0),
		make([]Middleware, 0),
	}
}

func (b *RoutesGroup) WithRoutes(routes ...*Route) *RoutesGroup {
	b.routes = append(b.routes, routes...)
	return b
}

func (b *RoutesGroup) WithGroups(builder ...*RoutesGroup) *RoutesGroup {
	b.routesGroups = append(b.routesGroups, builder...)
	return b
}

func (b *RoutesGroup) WithMiddlewares(middlewares ...Middleware) *RoutesGroup {
	b.middlewares = append(b.middlewares, middlewares...)
	return b
}

func (b *RoutesGroup) BuildRouter(headRouter *mux.Router) {
	subrouter := headRouter.PathPrefix(b.prefix).Subrouter()
	for _, route := range b.routes {
		subrouter.Handle(
			route.path,
			toHttpHandler(route.Middleware(b.middlewares...).handler),
		).Methods("OPTIONS", string(route.method))
	}

	if len(b.routesGroups) == 0 {
		return
	}

	for _, group := range b.routesGroups {
		group.middlewares = append(group.middlewares, b.middlewares...)
		group.BuildRouter(subrouter)
	}
}
