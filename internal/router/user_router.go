// Package router
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vcnt72/go-boilerplate/internal/handler"
)

func NewUserRouter(router *gin.Engine, userHandler *handler.UserHandler) {
	v1 := router.Group("v1")

	v1.POST("users", userHandler.Create())
}
