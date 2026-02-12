// Package config
package config

import (
	"os"

	"github.com/joho/godotenv"
)

type env struct {
	Port  string
	DBUrl string
}

var Env env

func Load() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	Env = env{
		Port:  os.Getenv("PORT"),
		DBUrl: os.Getenv("DB_URL"),
	}
}
