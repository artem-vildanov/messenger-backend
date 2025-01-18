package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	PgPassword string
	PgUser     string
	PgHost     string
	PgPort     string
	PgDb       string
	PgSSL      string
	RedisHost  string
	RedisPort  string
	AppHost    string
	AppPort    string
	Origin     string
	SessionTTL int // in minutes
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

	e.RedisHost = os.Getenv("REDIS_HOST")
	e.RedisPort = os.Getenv("REDIS_PORT")

	e.AppHost = os.Getenv("APP_HOST")
	e.AppPort = os.Getenv("APP_PORT")

	e.Origin = os.Getenv("ORIGIN")
	e.SessionTTL, err = strconv.Atoi(os.Getenv("SESSION_TTL")) // in minutes
	if err != nil {
		log.Panicf("Failed to convert SESSION_TTL to int: %s", err.Error())
	}
	return e
}
