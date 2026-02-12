package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vcnt72/go-boilerplate/internal/handler"
)

func New(router *gin.Engine, handlers handler.Handlers) {
	NewUserRouter(router, handlers.UserHandler)
	NewWalletRouter(router, handlers.WalletHandler)
}
