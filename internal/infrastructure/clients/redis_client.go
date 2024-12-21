package clients

import (
	"context"
	"messenger/internal/infrastructure/config"
	"messenger/internal/app/errors"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	connection *redis.Client
}

func (r *RedisClient) Construct(env *config.Env) {
	r.connection = redis.NewClient(&redis.Options{
		Addr: env.GetRedisAddr(),
	})
	if err := r.connection.Ping(context.Background()).Err(); err != nil {
		errors.ClientConnectionPanic("redis", err.Error())
	}
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.connection
}

func (r *RedisClient) CloseConnection() {
	r.connection.Close()
}
