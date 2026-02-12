// Package router
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vcnt72/go-boilerplate/internal/handler"
)

func NewWalletRouter(router *gin.Engine, walletHandler *handler.WalletHandler) {
	v1 := router.Group("v1")

	v1.GET("wallets/balance", walletHandler.GetBalance())
	v1.POST("wallets/withdraw", walletHandler.Withdraw())
}
