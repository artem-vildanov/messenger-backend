package inits

import (
	"context"
	"fmt"
	"log"
	"messenger/internal/infrastructure/config"
	"time"

	"github.com/go-redis/redis/v8"
)

func InitRedis(env *config.Env) (*redis.Client, func() error) {
	var err error
	conn := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			env.RedisHost,
			env.RedisPort,
		),
	})

	for retry := range env.RedisConnectRetries {
		err = conn.Ping(context.Background()).Err()
		if err == nil {
			break
		}
		log.Printf("redis connect retry %d...", retry)
		time.Sleep(time.Second)
	}

	if err != nil {
		log.Panicf("failed to connect to redis: %s", err.Error())
	}

	return conn, conn.Close
}
