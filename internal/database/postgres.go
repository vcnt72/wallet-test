package database

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/vcnt72/go-boilerplate/internal/config"
)

func NewPostgres() *sqlx.DB {
	db := sqlx.MustOpen("pgx", config.Env.DBUrl)

	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
