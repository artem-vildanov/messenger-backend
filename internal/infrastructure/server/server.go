package server

import (
	"log"
	"messenger/internal/app/middlewares"
	"messenger/internal/app/routes"
	"messenger/internal/infrastructure/clients"
	"messenger/internal/infrastructure/config"
	"messenger/internal/infrastructure/di"
	"messenger/internal/infrastructure/server/router"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	di *di.DependencyContainer
}

func New() *Server {
	server := &Server{
		mux.NewRouter(),
		di.NewDependencyContainer(),
	}
	return server
}

func (s *Server) Run() {
	s.initRouter()
	addr := di.FindDependency[config.Env](s.di).GetAppAddr()

	log.Printf("Server is running on %s...\n", addr)
	http.ListenAndServe(addr, s.router)
}

func (s *Server) initRouter() {
	router.RootGroup().WithGroups(
		routes.Api(s.di),
		routes.Ws(s.di),
	).WithMiddlewares(
		di.Provide[middlewares.CorsMiddleware](s.di),
	).BuildRouter(s.router)
}


func (s *Server) BeforeShutdown() {
	di.FindDependency[clients.PostgresClient](s.di).CloseConnection()
	// di.FindDependency[clients.RedisClient](s.di).CloseConnection()
}
