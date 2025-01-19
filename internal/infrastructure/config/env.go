package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	PgPassword          string
	PgUser              string
	PgHost              string
	PgPort              string
	PgDb                string
	PgSSL               string
	PgConnectRetries    int
	RedisHost           string
	RedisPort           string
	RedisConnectRetries int
	AppHost             string
	AppPort             string
	Origin              string
	SessionTTL          int // in minutes
}

func LoadEnv() *Env {
	e := &Env{}

	err := godotenv.Load()
	if err != nil {
		log.Panicf("Failed to load .env: %v", err)
	}

	e.PgPassword = os.Getenv("PG_PASSWORD")
	e.PgUser = os.Getenv("PG_USER")
	e.PgHost = os.Getenv("PG_HOST")
	e.PgPort = os.Getenv("PG_PORT")
	e.PgDb = os.Getenv("PG_DB")
	e.PgSSL = os.Getenv("PG_SSL")
	e.PgConnectRetries, err = strconv.Atoi(os.Getenv("PG_CONNECT_RETRIES"))
	if err != nil {
		convertIntPanic("PG_CONNECT_RETRIES", err)
	}

	e.RedisHost = os.Getenv("REDIS_HOST")
	e.RedisPort = os.Getenv("REDIS_PORT")
	e.RedisConnectRetries, err = strconv.Atoi(os.Getenv("REDIS_CONNECT_RETRIES"))
	if err != nil {
		convertIntPanic("REDIS_CONNECT_RETRIES", err)
	}

	e.AppHost = os.Getenv("APP_HOST")
	e.AppPort = os.Getenv("APP_PORT")

	e.Origin = os.Getenv("ORIGIN")
	e.SessionTTL, err = strconv.Atoi(os.Getenv("SESSION_TTL")) // in minutes
	if err != nil {
		convertIntPanic("SESSION_TTL", err)
	}
	return e
}

func convertIntPanic(envParam string, err error) {
	log.Panicf(
		"Failed to convert %s to int: %s",
		envParam,
		err.Error(),
	)
}
