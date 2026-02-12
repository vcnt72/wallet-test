// Package handler
package handler

import "github.com/vcnt72/go-boilerplate/internal/service"

type Handlers struct {
	UserHandler   *UserHandler
	WalletHandler *WalletHandler
}

func New(services service.Services) Handlers {
	return Handlers{
		UserHandler:   NewUserHandler(services.UserService),
		WalletHandler: NewWalletHandler(services.WalletService),
	}
}
