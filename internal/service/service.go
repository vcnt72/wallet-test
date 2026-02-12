package service

import "github.com/vcnt72/go-boilerplate/internal/repository"

type Services struct {
	UserService   *UserService
	WalletService *WalletService
}

func New(repositories repository.Repositories) Services {
	return Services{
		UserService: NewUserService(
			repositories.UserRepository,
			repositories.WalletRepository,
			repositories.LedgerRepository,
			repositories.TxProvider,
		),
		WalletService: NewWalletService(
			repositories.WalletRepository,
			repositories.LedgerRepository,
			repositories.TxProvider,
		),
	}
}
