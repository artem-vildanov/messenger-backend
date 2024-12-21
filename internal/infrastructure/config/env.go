package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	postgresPassword string
	postgresUser     string
	postgresHost     string
	postgresPort     string
	postgresDb       string

	redisHost string
	redisPort string

	appHost string
	appPort string

	sessionTTL int // in minutes
}

func (e *Env) Construct() {
	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Fatalf("Failed to load .env: %v", err)
	}

	e.postgresPassword = os.Getenv("POSTGRES_PASSWORD")
	e.postgresUser = os.Getenv("POSTGRES_USER")
	e.postgresHost = os.Getenv("POSTGRES_HOST")
	e.postgresPort = os.Getenv("POSTGRES_PORT")
	e.postgresDb = os.Getenv("POSTGRES_DB")

	e.redisHost = os.Getenv("REDIS_HOST")
	e.redisPort = os.Getenv("REDIS_PORT")

	e.appHost = os.Getenv("APP_HOST")
	e.appPort = os.Getenv("APP_PORT")

	e.sessionTTL, err = strconv.Atoi(os.Getenv("SESSION_TTL")) // in minutes
	if err != nil {
		log.Fatalf("Failed to convert SESSION_TTL to int: %s", err.Error())
	}
}

type PostgresAddr string

func (e *Env) GetPostgresAddr() PostgresAddr {
	return PostgresAddr(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			e.postgresUser,
			e.postgresPassword,
			e.postgresHost,
			e.postgresPort,
			e.postgresDb,
		),
	)
}

func (a PostgresAddr) WithSllDisabled() PostgresAddr {
	a += "?sslmode=disable"
	return a
}

func (a PostgresAddr) String() string {
	return string(a)
}

func (e *Env) GetRedisAddr() string {
	return fmt.Sprintf(
		"%s:%s", 
		e.redisHost, 
		e.redisPort,
	)
}

func (e *Env) GetAppAddr() string {
	return fmt.Sprintf(
		"%s:%s", 
		e.appHost, 
		e.appPort,
	)
}

func (e *Env) GetSessionTTL() int {
	return e.sessionTTL
}