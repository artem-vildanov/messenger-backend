package bootstrap

import (
	"fmt"
	"log"
	"messenger/internal/domain/services"
	"messenger/internal/domain/storage"
	"messenger/internal/infrastructure/config"
	"messenger/internal/infrastructure/inits"
	"messenger/internal/infrastructure/pubsub"
	"messenger/internal/infrastructure/utils/router_utils"
	"messenger/internal/presentation/handlers"
	"messenger/internal/presentation/middlewares"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type StorageRegistry struct {
	PgConn *sqlx.DB
	ChatStorage    *storage.ChatStorage
	MessageStorage *storage.MessageStorage
	SessionStorage *storage.SessionStorage
	UserStorage    *storage.UserStorage
}

func NewStorageRegistry(pg *sqlx.DB) *StorageRegistry {
	// init storages
	return &StorageRegistry{
		pg,
		storage.NewChatStorage(pg),
		storage.NewMessageStorage(pg),
		storage.NewSessionStorage(pg),
		storage.NewUserStorage(pg),
	}
}

type ServiceRegistry struct {
	ChatService     *services.ChatService
	SesssionService *services.SessionService
}

func NewServiceRegistry(
	env *config.Env,
	redis *redis.Client,
	storageRegistry *StorageRegistry,
) *ServiceRegistry {
	return &ServiceRegistry{
		services.NewChatService(
			pubsub.NewPubsubClient(redis),
			storageRegistry.MessageStorage,
		),
		services.NewSessionService(
			storageRegistry.SessionStorage,
			storageRegistry.UserStorage,
			env,
		),
	}
}

type HandlerRegistry struct {
	AuthHandler *handlers.AuthHandler
	ChatHandler *handlers.ChatHandler
	UserHandler *handlers.UserHandler
}

func NewHandlerRegistry(
	storageRegistry *StorageRegistry,
	serviceRegistry *ServiceRegistry,
	env *config.Env,
) *HandlerRegistry {
	return &HandlerRegistry{
		handlers.NewAuthHandler(
			storageRegistry.SessionStorage,
			serviceRegistry.SesssionService,
			env,
		),
		handlers.NewChatHandler(
			storageRegistry.ChatStorage,
			storageRegistry.MessageStorage,
			serviceRegistry.ChatService,
		),
		handlers.NewUserHandler(
			storageRegistry.UserStorage,
		),
	}
}

type MiddlewareRegistry struct {
	AuthMiddleware    router_utils.Middleware
	LoggingMiddleware router_utils.Middleware
	CorsMiddleware    router_utils.Middleware
}

func NewMiddlewareRegistry(
	serviceRegistry *ServiceRegistry,
	env *config.Env,
) *MiddlewareRegistry {
	return &MiddlewareRegistry{
		middlewares.NewAuthMiddleware(serviceRegistry.SesssionService),
		middlewares.NewLogginMiddleware(),
		middlewares.NewCorsMiddleware(env),
	}
}

type App struct {
	Env                *config.Env
	StorageRegistry *StorageRegistry
	ServiceRegistry *ServiceRegistry
	HandlerRegistry    *HandlerRegistry
	MiddlewareRegistry *MiddlewareRegistry
	Cleanup            func()
}

func NewApp() *App {
	env := config.LoadEnv()

	// init connections
	pg, closePgConn := inits.InitPostgres(env)
	redis, closeRedisConn := inits.InitRedis(env)

	storageRegistry := NewStorageRegistry(pg)
	serviceRegistry := NewServiceRegistry(env, redis, storageRegistry)
	handlerRegistry := NewHandlerRegistry(storageRegistry, serviceRegistry, env)
	middlewareRegistry := NewMiddlewareRegistry(serviceRegistry, env)

	return &App{
		env,
		storageRegistry,
		serviceRegistry,
		handlerRegistry,
		middlewareRegistry,
		func() {
			closePgConn()
			closeRedisConn()
		},
	}
}

func (s *App) Run(routes *router_utils.RoutesGroup) {
	router := mux.NewRouter()
	routes.BuildRouter(router)

	addr := fmt.Sprintf(
		"%s:%s",
		s.Env.AppHost,
		s.Env.AppPort,
	)

	log.Printf("Server is running on %s...\n", addr)
	http.ListenAndServe(addr, router)
}
