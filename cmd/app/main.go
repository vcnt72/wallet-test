package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vcnt72/go-boilerplate/internal/config"
	"github.com/vcnt72/go-boilerplate/internal/database"
	"github.com/vcnt72/go-boilerplate/internal/handler"
	"github.com/vcnt72/go-boilerplate/internal/repository"
	"github.com/vcnt72/go-boilerplate/internal/router"
	"github.com/vcnt72/go-boilerplate/internal/service"
)

func main() {
	config.Load()
	db := database.NewPostgres()
	routerEngine := gin.Default()

	repositories := repository.New(db)

	services := service.New(repositories)

	handlers := handler.New(services)

	router.New(routerEngine, handlers)

	if err := routerEngine.Run(fmt.Sprintf(":%s", config.Env.Port)); err != nil {
		log.Fatalf("%v", err)
	}
}
