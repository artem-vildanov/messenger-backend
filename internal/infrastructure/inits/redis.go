package inits

import (
	"context"
	"fmt"
	"log"
	"messenger/internal/infrastructure/config"

	"github.com/go-redis/redis/v8"
)

func InitRedis(env *config.Env) (*redis.Client, func() error) {
	conn := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			env.RedisHost,
			env.RedisPort,
		),
	})
	if err := conn.Ping(context.Background()).Err(); err != nil {
		log.Panicf("failed to connect to redis: %s", err.Error())
	}
	return conn, conn.Close
}
